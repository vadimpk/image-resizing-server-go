package main

import (
	"github.com/vadimpk/image-resizing-server-go/internal/api/delivery/http"
	"github.com/vadimpk/image-resizing-server-go/internal/api/server"
	"github.com/vadimpk/image-resizing-server-go/internal/api/service"
	"log"
)

func main() {
	services := service.NewServices()

	handler := http.NewHandler(services)
	r := handler.Init()

	srv := server.NewServer(r)

	if err := srv.Run(); err != nil {
		log.Fatalf("Error while running server: %s", err.Error())
	}

}
