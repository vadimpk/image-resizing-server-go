package service

import (
	"github.com/vadimpk/image-resizing-server-go/internal/api/publisher"
	"github.com/vadimpk/image-resizing-server-go/internal/api/repository"
)

type Services struct {
	Uploading
	Downloading
}

type Uploading interface {
	Upload(file []byte) (string, error)
}

type Downloading interface {
	Download(ID string) ([]byte, error) // TODO: return image instead of bytes
}

func NewServices(publisher publisher.Publisher, repository *repository.Repository) *Services {
	return &Services{
		Uploading:   NewUploadingService(publisher),
		Downloading: NewDownloadingService(repository),
	}
}
