package main

import (
	"github.com/vadimpk/image-resizing-server-go/internal/api/delivery/http"
	"github.com/vadimpk/image-resizing-server-go/internal/api/publisher"
	"github.com/vadimpk/image-resizing-server-go/internal/api/repository"
	"github.com/vadimpk/image-resizing-server-go/internal/api/server"
	"github.com/vadimpk/image-resizing-server-go/internal/api/service"
	"github.com/vadimpk/image-resizing-server-go/pkg/queue/rabbitmq"
	"github.com/vadimpk/image-resizing-server-go/pkg/queue/rabbitmq/producer"
	"log"
)

func main() {

	rabbit, err := rabbitmq.NewRabbitMQConn("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	var p publisher.Publisher
	p, err = producer.NewProducer(rabbit, producer.QueueConfig{
		Name:         "test",
		Durable:      false,
		DeleteUnused: false,
		Exclusive:    false,
		NoWait:       false,
	})
	if err != nil {
		log.Fatal(err)
	}

	services := service.NewServices(p, repository.NewRepository())

	handler := http.NewHandler(services)
	r := handler.Init()

	srv := server.NewServer(r)

	if err := srv.Run(); err != nil {
		log.Fatalf("Error while running server: %s", err.Error())
	}
}
