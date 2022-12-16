package service

type UploadingService struct {
	// publisher
}

func NewUploadingService() *UploadingService {
	return &UploadingService{}
}

func (s *UploadingService) Upload(file []byte) (string, error) {

	go func() {
		// publish
	}()

	// TODO: generate id

	return "id", nil
}
