package rabbitmq

import (
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

const (
	_retryTimes     = 5
	_backOffSeconds = 2
)

func NewRabbitMQConn(rabbitMqURL string) (*amqp.Connection, error) {

	var amqpConn *amqp.Connection
	count := 0
	for {
		connection, err := amqp.Dial(rabbitMqURL)
		if err == nil {
			amqpConn = connection
			break
		}

		log.Printf("Couldn't connect to RabbitMQ at %s...\n\n", rabbitMqURL)
		count++

		if count > _retryTimes {
			return nil, errors.New("couldn't connect to RabbitMQ")
		}

		log.Printf("Retrying in %d seconds ...", _backOffSeconds)
		time.Sleep(_backOffSeconds * time.Second)
	}

	log.Println("Connected to RabbitMQ!")

	return amqpConn, nil
}
