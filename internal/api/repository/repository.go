package repository

import "context"

type ImgRepository interface {
	GetImage(id string, resolution int) ([]byte, error)
	Close(ctx context.Context) chan struct{}
}

type Repository struct {
	ImgRepository
}

func NewRepository() *Repository {
	return &Repository{&TempImgRepo{}}
}

type TempImgRepo struct {
}

func (r *TempImgRepo) GetImage(id string, resolution int) ([]byte, error) {
	return nil, nil
}

func (r *TempImgRepo) Close(ctx context.Context) chan struct{} {
	return nil
}
