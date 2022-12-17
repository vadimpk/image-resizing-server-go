package main

import (
	"context"
	"github.com/vadimpk/image-resizing-server-go/internal/api/delivery/http"
	"github.com/vadimpk/image-resizing-server-go/internal/api/publisher"
	"github.com/vadimpk/image-resizing-server-go/internal/api/repository"
	"github.com/vadimpk/image-resizing-server-go/internal/api/server"
	"github.com/vadimpk/image-resizing-server-go/internal/api/service"
	"github.com/vadimpk/image-resizing-server-go/pkg/queue/rabbitmq"
	"github.com/vadimpk/image-resizing-server-go/pkg/queue/rabbitmq/producer"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	rabbit := rabbitmq.NewRabbitMQ()
	err := rabbit.Connect("amqp://guest:guest@localhost:5672/")
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

	_, cancel := context.WithCancel(context.Background())

	srv.Run()

	defer shutdown(cancel, srv, p)

	waitShutdown()
}

func shutdown(cancel context.CancelFunc, srv *server.Server, p publisher.Publisher) {
	cancel()
	ctx, cancelTimeout := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelTimeout()

	doneHTTP := srv.Stop(ctx)
	doneRabbit := p.Close(ctx)

	waitUntilIsDoneOrCanceled(ctx, doneHTTP, doneRabbit)
	time.Sleep(time.Millisecond * 200)
}

func waitShutdown() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sigc
	log.Printf("signal [%v] received, canceling everything...\n", s)
}

func waitUntilIsDoneOrCanceled(ctx context.Context, dones ...chan struct{}) {
	done := make(chan struct{})
	go func() {
		for _, d := range dones {
			<-d
		}
		close(done)
	}()
	select {
	case <-done:
		log.Println("all done")
	case <-ctx.Done():
		log.Println("canceled")
	}
}
