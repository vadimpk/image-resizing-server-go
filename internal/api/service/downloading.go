package service

type DownloadingService struct {
	// repo
}

func NewDownloadingService() *DownloadingService {
	return &DownloadingService{}
}

func (s *DownloadingService) Download(ID string) ([]byte, error) {
	return nil, nil
}
