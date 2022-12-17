package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
	"time"
)

const (
	_retryTimes     = 5
	_backOffSeconds = 2
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Wg      *sync.WaitGroup
}

func NewRabbitMQ() *RabbitMQ {
	return &RabbitMQ{
		Wg: &sync.WaitGroup{},
	}
}

func (r *RabbitMQ) Connect(rabbitMqURL string) error {

	var conn *amqp.Connection
	count := 0
	for {
		connection, err := amqp.Dial(rabbitMqURL)
		if err == nil {
			conn = connection
			break
		}

		log.Printf("Couldn't connect to RabbitMQ at %s...\n\n", rabbitMqURL)
		count++

		if count > _retryTimes {
			return err
		}

		log.Printf("Retrying in %d seconds ...", _backOffSeconds)
		time.Sleep(_backOffSeconds * time.Second)
	}

	log.Println("Connected to RabbitMQ!")

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Couldn't open RabbitMQ channel: [%s]\n", err)
		return err
	}

	r.Conn = conn
	r.Channel = ch
	return nil
}

// Close gracefully closes the rabbitMQ connection, trying to finish all the jobs first
// Returns done channel that notifies when all connections are closed
func (r *RabbitMQ) Close(ctx context.Context) chan struct{} {
	done := make(chan struct{})

	doneWaiting := make(chan struct{})
	go func() {
		r.Wg.Wait()
		close(doneWaiting)
	}()

	go func() {
		defer close(done)
		select { // either waits for the messages to process or timeout from context
		case <-doneWaiting:
		case <-ctx.Done():
		}
		r.closeConnections()
	}()

	return done
}

func (r *RabbitMQ) closeConnections() {
	if r.Channel != nil {
		err := r.Channel.Close()
		if err != nil {
			log.Printf("error closing channel: [%s]\n", err)
		}
	}

	if r.Conn != nil {
		err := r.Conn.Close()
		if err != nil {
			log.Printf("error closing connection: [%s]\n", err)
		}
	}
}
