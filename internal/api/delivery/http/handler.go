package http

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vadimpk/image-resizing-server-go/internal/api/service"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	r.HandleFunc("/download/{id}", h.HandleDownload)

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

func (h *Handler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	quality := r.URL.Query().Get("quality")

	if quality == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resolution, err := strconv.Atoi(quality)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	d, filename, err := h.service.Download(id, resolution)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ct, err := determineContentType(filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ct)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	io.Copy(w, bytes.NewReader(d))
}

func determineContentType(filename string) (string, error) {
	parts := strings.Split(filename, ".")
	if len(parts) != 2 {
		return "", errors.New("couldn't determine content type")
	}
	switch parts[1] {
	case "jpeg", "jpg":
		return "image/jpeg", nil
	case "png":
		return "image/png", nil
	default:
		return "", errors.New("couldn't determine content type")
	}
}
