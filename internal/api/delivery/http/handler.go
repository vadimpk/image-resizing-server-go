package http

import (
	"github.com/gorilla/mux"
	"github.com/vadimpk/image-resizing-server-go/internal/api/service"
	"io/ioutil"
	"log"
	"net/http"
)

type Handler struct {
	service *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services}
}

func (h *Handler) Init() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/upload", h.HandleUpload)

	return r
}

func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	log.Println("Receiving file...")

	file, headers, err := r.FormFile("myFile")
	if err != nil {
		log.Printf("Error receiving file: [%s]", err)
		return
	}

	log.Printf("Received File: %+v\n", headers.Filename)
	log.Printf("File Size: %+v\n", headers.Size)
	log.Printf("MIME Header: %+v\n", headers.Header)

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file: [%s]", err)
		return
	}

	id, err := h.service.Upload(bytes)
	if err != nil {
		log.Printf("Error uploading file: [%s]", err)
		return
	}

	_, err = w.Write([]byte(id))
	if err != nil {
		log.Printf("Error writing response: [%s]", err)
		return
	}
}
