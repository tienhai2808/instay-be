package port

import "context"

type CtxTxKey struct{}

type TxProvider interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}
