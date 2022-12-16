package service

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

func NewServices() *Services {
	return &Services{
		Uploading:   NewUploadingService(),
		Downloading: NewDownloadingService(),
	}
}
