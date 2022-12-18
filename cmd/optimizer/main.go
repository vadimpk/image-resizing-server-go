package main

import (
	"context"
	"github.com/vadimpk/image-resizing-server-go/internal/optimizer/consumer"
	"github.com/vadimpk/image-resizing-server-go/internal/optimizer/repository"
	"github.com/vadimpk/image-resizing-server-go/internal/optimizer/services"
	"github.com/vadimpk/image-resizing-server-go/pkg/queue/rabbitmq"
	rabbitconsumer "github.com/vadimpk/image-resizing-server-go/pkg/queue/rabbitmq/consumer"
	"github.com/vadimpk/image-resizing-server-go/pkg/storage/filestorage"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	var repo repository.Repository
	repo = filestorage.NewStorage()

	var service services.Service
	service = services.NewOptimizer(repo)

	rabbit := rabbitmq.NewRabbitMQ()
	err := rabbit.Connect("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	var c consumer.Consumer
	c, err = rabbitconsumer.NewConsumer(rabbit)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err = c.Consume(ctx, service.Optimize)
		if err != nil {
			log.Fatalf("couldn't consume messages: [%s]\n", err)
		}
	}()

	defer shutdown(cancel, c, repo)

	waitShutdown()
}

func shutdown(cancel context.CancelFunc, c consumer.Consumer, r repository.Repository) {
	cancel()
	ctx, cancelTimeout := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelTimeout()

	doneRabbit := c.Close(ctx)
	doneRepo := r.Close(ctx)

	waitUntilIsDoneOrCanceled(ctx, doneRabbit, doneRepo)
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
