package service

import (
	"context"
	"errors"
	"github.com/teris-io/shortid"
	http2 "github.com/vadimpk/image-resizing-server-go/internal/api/delivery/http"
	"github.com/vadimpk/image-resizing-server-go/internal/api/publisher"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

type UploadingService struct {
	publisher publisher.Publisher
	sid       *shortid.Shortid
}

func NewUploadingService(publisher publisher.Publisher, sid *shortid.Shortid) *UploadingService {
	return &UploadingService{publisher, sid}
}

func (s *UploadingService) Upload(file multipart.File) (string, error) {

	// reading file
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("error reading file: [%s]", err)
		return "", err
	}

	// determining content type
	contentType := http.DetectContentType(bytes)
	switch contentType {
	case "image/jpeg", "image/png":
		break
	default:
		return "", http2.ErrInvalidContentType
	}

	// generating id
	id, err := s.sid.Generate()
	if err != nil {
		return "", errors.New("couldn't generate id")
	}

	go func() {
		h := map[string]interface{}{
			"img-type": contentType,
			"id":       id,
		}
		err = s.publisher.Publish(context.Background(), bytes, h)
		if err != nil {
			log.Printf("couldn't publish image: [%s]\n", err)
		}
		//TODO: if error - add id to special list, and when request comes for download, say there was some error uploading
	}()

	return id, nil
}
