package repository

import "context"

type SchemeRepository interface {
	FindAll(ctx context.Context) ([]string, error)
}
