package service

import "context"

// TxManager — абстракция над транзакционным менеджером.
// Сервисы знают только про этот интерфейс.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
