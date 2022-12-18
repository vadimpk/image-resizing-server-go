package main

import (
	"context"
	"github.com/vadimpk/image-resizing-server-go/internal/optimizer/config"
	"github.com/vadimpk/image-resizing-server-go/internal/optimizer/consumer"
	"github.com/vadimpk/image-resizing-server-go/internal/optimizer/repository"
	"github.com/vadimpk/image-resizing-server-go/internal/optimizer/service"
	"github.com/vadimpk/image-resizing-server-go/pkg/queue/rabbitmq"
	rabbitconsumer "github.com/vadimpk/image-resizing-server-go/pkg/queue/rabbitmq/consumer"
	"github.com/vadimpk/image-resizing-server-go/pkg/storage/filestorage"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	repo     repository.Repository
	services service.Service
	c        consumer.Consumer
)

func main() {

	cfg, err := config.Init("configs/main")
	if err != nil {
		log.Fatalf("couldn't parse config: [%s]\n", err)
	}

	rabbit := rabbitmq.NewRabbitMQ()
	err = rabbit.Connect(cfg.Rabbit.URL + cfg.Rabbit.Port + "/")
	if err != nil {
		log.Fatalf("couldn't connect to RabbitMQ: []%s\n", err)
	}

	c, err = rabbitconsumer.NewConsumer(rabbit)
	if err != nil {
		log.Fatalf("couldn't initialize queue consumer: [%s]\n", err)
	}

	repo = filestorage.NewStorage(cfg.FileStorage.DirPath)
	services = service.NewOptimizer(repo)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err = c.Consume(ctx, services.Optimize)
		if err != nil {
			log.Fatalf("couldn't consume messages: [%s]\n", err)
		}
	}()

	defer shutdown(cancel)
	waitShutdown()
}

// shutdown gracefully stops all services after given timeout
func shutdown(cancel context.CancelFunc) {
	cancel()
	ctx, cancelTimeout := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelTimeout()

	doneRabbit := c.Close(ctx)
	doneRepo := repo.Close(ctx)

	waitUntilIsDoneOrCanceled(ctx, doneRabbit, doneRepo)
	time.Sleep(time.Millisecond * 200)
}

// waitShutdown waits for stop signal to stop the program
func waitShutdown() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sigc
	log.Printf("signal [%v] received, canceling everything...\n", s)
}

// waitUntilIsDoneOrCanceled implements graceful shutdown, it waits for all the processes to be finished
// or for the context to timeout
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
		log.Println("All processes finished")
	case <-ctx.Done():
		log.Println("Some processes were killed after delay")
	}
}
