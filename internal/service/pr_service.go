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
	var merged *domain.PullRequest

	err := s.tx.Do(ctx, func(txCtx context.Context) error {
		// 1. достаём PR
		pr, err := s.pr.GetByID(txCtx, prID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return domain.NewNotFoundError("pull_request")
			}
			return err
		}

		// 2. если уже MERGED — просто возвращаем как есть (идемпотентность)
		if pr.Status == domain.PRStatusMerged {
			merged = pr
			return nil
		}

		// 3. помечаем как MERGED и проставляем mergedAt
		now := time.Now().UTC()
		pr.Status = domain.PRStatusMerged
		pr.MergedAt = &now

		if err := s.pr.Update(txCtx, *pr); err != nil {
			return err
		}

		merged = pr
		return nil
	})

	if err != nil {
		return nil, err
	}

	return merged, nil
}

func (s *prService) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*domain.PullRequest, string, error) {
	var updated *domain.PullRequest
	var replacedBy string

	err := s.tx.Do(ctx, func(txCtx context.Context) error {
		// 1. достаём PR
		pr, err := s.pr.GetByID(txCtx, prID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return domain.NewNotFoundError("pull_request")
			}
			return err
		}

		// 2. нельзя менять ревьюверов после MERGED
		if pr.Status == domain.PRStatusMerged {
			return domain.NewPRMergedError(prID)
		}

		// 3. проверяем, что старый ревьювер вообще назначен
		idx := -1
		for i, rid := range pr.AssignedReviewers {
			if rid == oldReviewerID {
				idx = i
				break
			}
		}
		if idx == -1 {
			return domain.NewNotAssignedError(oldReviewerID, prID)
		}

		// 4. убеждаемся, что такой пользователь существует
		oldReviewer, err := s.users.GetByID(txCtx, oldReviewerID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return domain.NewNotFoundError("user")
			}
			return err
		}

		// 5. берём активных участников команды старого ревьювера
		candidates, err := s.users.GetActiveTeamMembers(txCtx, oldReviewer.TeamName)
		if err != nil {
			return err
		}

		// 6. фильтруем: убираем самого старого ревьювера и уже назначенных ревьюверов,
		// чтобы не дублировать одного и того же человека
		current := make(map[string]struct{}, len(pr.AssignedReviewers))
		for _, rid := range pr.AssignedReviewers {
			current[rid] = struct{}{}
		}

		filtered := make([]domain.User, 0, len(candidates))
		for _, u := range candidates {
			if u.ID == oldReviewerID {
				continue
			}
			// не назначаем второго такого же ревьювера
			if _, exists := current[u.ID]; exists {
				continue
			}
			filtered = append(filtered, u)
		}

		if len(filtered) == 0 {
			return domain.NewNoCandidateError(oldReviewer.TeamName)
		}

		// 7. выбираем одного случайного кандидата
		newReviewerID := pickRandomReviewers(filtered, 1)[0]

		// 8. заменяем в списке ревьюверов
		pr.AssignedReviewers[idx] = newReviewerID

		// 9. сохраняем изменения
		if err := s.pr.Update(txCtx, *pr); err != nil {
			return err
		}

		updated = pr
		replacedBy = newReviewerID
		return nil
	})

	if err != nil {
		return nil, "", err
	}

	return updated, replacedBy, nil
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
