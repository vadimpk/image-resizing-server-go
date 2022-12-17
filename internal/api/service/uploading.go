package service

import (
	"context"
	"github.com/vadimpk/image-resizing-server-go/internal/api/publisher"
)

type UploadingService struct {
	publisher publisher.Publisher
}

func NewUploadingService(queue publisher.Publisher) *UploadingService {
	return &UploadingService{queue}
}

func (s *UploadingService) Upload(file []byte) (string, error) {

	// TODO: generate id
	id := "id"

	go func() {
		headers := map[string]interface{}{
			"img-type": "jpeg",
			"id":       id,
		}
		s.publisher.PublishImage(context.Background(), file, headers)
	}()

	return id, nil
}
