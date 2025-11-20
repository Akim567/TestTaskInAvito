package tx

import (
	"context"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	trmManager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
)

// Manager — обёртка над go-transaction-manager, реализующая наш интерфейс service.TxManager.
type Manager struct {
	m *trmManager.Manager
}

// NewManager создаёт транзакционный менеджер для sqlx.DB.
func NewManager(db *sqlx.DB) *Manager {
	trFactory := trmsqlx.NewDefaultFactory(db)
	m := trmManager.Must(trFactory)

	return &Manager{m: m}
}

// Do реализует интерфейс TxManager
func (tm *Manager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.m.Do(ctx, fn)
}

// NewCtxGetter нужен будет репозиториям, чтобы доставать из контекста либо транзакцию, либо обычный db.
func NewCtxGetter() *trmsqlx.CtxGetter {
	return trmsqlx.DefaultCtxGetter
}
