package service

import "github.com/vadimpk/image-resizing-server-go/internal/api/repository"

type DownloadingService struct {
	repo *repository.Repository
}

func NewDownloadingService(repository *repository.Repository) *DownloadingService {
	return &DownloadingService{repository}
}

func (s *DownloadingService) Download(ID string) ([]byte, error) {
	return nil, nil
}
