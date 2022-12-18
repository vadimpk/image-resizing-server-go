package repository

import "context"

type Repository interface {
	Get(id string, resolution int) ([]byte, string, error)
	Close(ctx context.Context) chan struct{}
}
