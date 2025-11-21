package service

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"TestTaskInAvito/internal/domain"
	"TestTaskInAvito/internal/repo"
)

// PRService — управление PR'ами.
type PRService interface {
	CreatePR(ctx context.Context, pr domain.PullRequest) (*domain.PullRequest, error)
	MergePR(ctx context.Context, prID string) (*domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*domain.PullRequest, string, error)
}

type prService struct {
	pr    repo.PRRepository
	users repo.UserRepository
	teams repo.TeamRepository
	tx    TxManager
}

func NewPRService(
	pr repo.PRRepository,
	users repo.UserRepository,
	teams repo.TeamRepository,
	tx TxManager,
) PRService {
	return &prService{
		pr:    pr,
		users: users,
		teams: teams,
		tx:    tx,
	}
}

// ===== реализация интерфейса =====

// CreatePR создаёт новый PR, назначая до двух случайных активных ревьюверов из команды автора.
func (s *prService) CreatePR(ctx context.Context, in domain.PullRequest) (*domain.PullRequest, error) {
	var created *domain.PullRequest

	err := s.tx.Do(ctx, func(txCtx context.Context) error {
		// Проверяем, что PR с таким id ещё не существует
		existing, err := s.pr.GetByID(txCtx, in.ID)
		if err == nil && existing != nil {
			return domain.NewPRExistsError(in.ID)
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			// неизвестная ошибка из репозитория
			return err
		}

		// Получаем автора
		author, err := s.users.GetByID(txCtx, in.Author)
		if err != nil {
			return err
		}

		// Берём активных участников команды автора
		candidates, err := s.users.GetActiveTeamMembers(txCtx, author.TeamName)
		if err != nil {
			return err
		}

		// Исключаем автора из списка
		filtered := make([]domain.User, 0, len(candidates))
		for _, u := range candidates {
			if u.ID == author.ID {
				continue
			}
			filtered = append(filtered, u)
		}

		// Выбираем до двух случайных ревьюверов
		reviewerIDs := pickRandomReviewers(filtered, 2)

		// Собираем доменный объект PR
		now := time.Now().UTC()
		newPR := domain.PullRequest{
			ID:                in.ID,
			Name:              in.Name,
			Author:            in.Author,
			Status:            domain.PRStatusOpen,
			AssignedReviewers: reviewerIDs,
			CreatedAt:         now,
			MergedAt:          nil,
		}

		// Сохраняем в БД
		if err := s.pr.Create(txCtx, newPR); err != nil {
			return err
		}

		created = &newPR
		return nil
	})

	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *prService) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	// TODO: реализовать идемпотентный merge
	panic("not implemented")
}

func (s *prService) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*domain.PullRequest, string, error) {
	// TODO: реализовать правила переназначения
	panic("not implemented")
}

// ===== хелперы =====

// единый rand для всего пакета (чтобы не переинициализировать каждый вызов)
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// pickRandomReviewers выбирает до count случайных пользователей и возвращает их ID.
func pickRandomReviewers(users []domain.User, count int) []string {
	if len(users) == 0 || count <= 0 {
		return nil
	}
	if len(users) <= count {
		ids := make([]string, 0, len(users))
		for _, u := range users {
			ids = append(ids, u.ID)
		}
		return ids
	}

	// копируем, чтобы не перемешивать исходный слайс
	tmp := make([]domain.User, len(users))
	copy(tmp, users)

	rnd.Shuffle(len(tmp), func(i, j int) {
		tmp[i], tmp[j] = tmp[j], tmp[i]
	})

	tmp = tmp[:count]
	ids := make([]string, 0, len(tmp))
	for _, u := range tmp {
		ids = append(ids, u.ID)
	}
	return ids
}
