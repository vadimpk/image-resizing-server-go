package service

import (
	"github.com/teris-io/shortid"
	"github.com/vadimpk/image-resizing-server-go/internal/api/publisher"
	"github.com/vadimpk/image-resizing-server-go/internal/api/repository"
	"mime/multipart"
)

type Services struct {
	Uploading
	Downloading
}

type Uploading interface {
	Upload(file multipart.File) (string, error)
}

type Downloading interface {
	Download(id string) ([]byte, error) // TODO: return image instead of bytes
}

func NewServices(publisher publisher.Publisher, repository *repository.Repository, sid *shortid.Shortid) *Services {
	return &Services{
		Uploading:   NewUploadingService(publisher, sid),
		Downloading: NewDownloadingService(repository),
	}
}
