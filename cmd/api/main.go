package main

import (
	"context"
	"github.com/teris-io/shortid"
	"github.com/vadimpk/image-resizing-server-go/internal/api/config"
	"github.com/vadimpk/image-resizing-server-go/internal/api/delivery/http"
	"github.com/vadimpk/image-resizing-server-go/internal/api/publisher"
	"github.com/vadimpk/image-resizing-server-go/internal/api/repository"
	"github.com/vadimpk/image-resizing-server-go/internal/api/server"
	"github.com/vadimpk/image-resizing-server-go/internal/api/service"
	"github.com/vadimpk/image-resizing-server-go/pkg/queue/rabbitmq"
	"github.com/vadimpk/image-resizing-server-go/pkg/queue/rabbitmq/producer"
	"github.com/vadimpk/image-resizing-server-go/pkg/storage/filestorage"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	repo     repository.Repository
	services *service.Services
	handler  *http.Handler
	srv      *server.Server
	pub      publisher.Publisher
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

	pub, err = producer.NewProducer(rabbit)
	if err != nil {
		log.Fatalf("couldn't initialize queue producer: [%s]\n", err)
	}

	sid, err := shortid.New(1, shortid.DefaultABC, 4512)
	if err != nil {
		log.Fatalf("couldn't create random id generator: [%s]\n", err)
	}

	repo = filestorage.NewStorage(cfg.FileStorage.DirPath)
	services = service.NewServices(pub, repo, sid)

	handler = http.NewHandler(services, cfg.Server.MaxFileSizeMB)
	r := handler.Init()

	srv = server.NewServer(cfg, r)

	_, cancel := context.WithCancel(context.Background())

	srv.Run()
	defer shutdown(cancel, cfg.Main.Timeout)
	waitShutdown()
}

// shutdown gracefully stops all services after given timeout
func shutdown(cancel context.CancelFunc, timeout time.Duration) {
	cancel()
	ctx, cancelTimeout := context.WithTimeout(context.Background(), timeout)
	defer cancelTimeout()

	doneHTTP := srv.Stop(ctx)
	doneRabbit := pub.Close(ctx)

	waitUntilIsDoneOrCanceled(ctx, doneHTTP, doneRabbit)
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
