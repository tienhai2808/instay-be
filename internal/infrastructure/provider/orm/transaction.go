package orm

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"gorm.io/gorm"
)

type txProviderImpl struct {
	db *gorm.DB
}

func NewTxProvider(db *gorm.DB) port.TxProvider {
	return &txProviderImpl{db}
}

func (p *txProviderImpl) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctxWithTx := context.WithValue(ctx, port.CtxTxKey{}, tx)
		return fn(ctxWithTx)
	})
}
