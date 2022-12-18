package service

import (
	"context"
	"errors"
	"github.com/teris-io/shortid"
	_ "github.com/teris-io/shortid"
	"github.com/vadimpk/image-resizing-server-go/internal/api/publisher"
	"io/ioutil"
	"log"
	"mime/multipart"
)

type UploadingService struct {
	publisher publisher.Publisher
	sid       *shortid.Shortid
}

func NewUploadingService(publisher publisher.Publisher, sid *shortid.Shortid) *UploadingService {
	return &UploadingService{publisher, sid}
}

func (s *UploadingService) Upload(file multipart.File, headers *multipart.FileHeader) (string, error) {

	//TODO: check image type compatability with service (jpeg/png)

	id, err := s.sid.Generate()
	if err != nil {
		return "", errors.New("couldn't generate id")
	}

	go func() {
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Printf("Error reading file: [%s]", err)
			return
		}

		h := map[string]interface{}{
			"img-type": headers.Header.Get("Content-Type"),
			"id":       id,
		}
		err = s.publisher.Publish(context.Background(), bytes, h)
		if err != nil {
			log.Printf("couldn't publish image: [%s]\n", err)
		}
		//TODO: if error - add id to special list, and when request comes for download, say there was some error
	}()

	return id, nil
}
