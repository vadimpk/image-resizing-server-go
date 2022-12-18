package repository

import (
	"context"
	"image"
)

type Repository interface {
	Save(id string, image image.Image, resolution int, imgType string) error
	Close(ctx context.Context) chan struct{}
}
