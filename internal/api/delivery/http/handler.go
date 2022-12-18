package http

import (
	"github.com/gorilla/mux"
	"github.com/vadimpk/image-resizing-server-go/internal/api/service"
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

	err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("Receiving file...")

	file, headers, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error receiving file: [%s]", err)
		return
	}

	log.Printf("Received File: %+v\n", headers.Filename)
	log.Printf("File Size: %+v\n", headers.Size)
	log.Printf("MIME Header: %+v\n", headers.Header)

	id, err := h.service.Upload(file)
	if err != nil {
		log.Printf("Error uploading file: [%s]", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(id))
	if err != nil {
		log.Printf("Error writing response: [%s]", err)
		w.WriteHeader(http.StatusOK)
		return
	}
}
