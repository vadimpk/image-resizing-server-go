package producer

import (
	"context"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/vadimpk/image-resizing-server-go/pkg/queue/rabbitmq"
	"log"
	"time"
)

type Producer struct {
	rabbit     *rabbitmq.RabbitMQ
	routingKey string
}

func NewProducer(rabbit *rabbitmq.RabbitMQ, queueConfig QueueConfig) (*Producer, error) {

	// TODO: try to create new
	if rabbit.Channel == nil {
		return nil, errors.New("RabbitMQ channel not found when creating producer")
	}

	q, err := rabbit.Channel.QueueDeclare(
		queueConfig.Name,
		queueConfig.Durable,
		queueConfig.DeleteUnused,
		queueConfig.Exclusive,
		queueConfig.NoWait,
		nil,
	)
	if err != nil {
		log.Printf("Couldn't create RabbitMQ publisher: [%s]\n", err)
		return nil, err
	}

	p := &Producer{
		rabbit:     rabbit,
		routingKey: q.Name,
	}

	return p, nil
}

func (p *Producer) Close(ctx context.Context) chan struct{} {
	return p.rabbit.Close(ctx)
}

func (p *Producer) Publish(ctx context.Context, body []byte, headers map[string]interface{}) error {

	// TODO: create new channel (reconnect)
	if p.rabbit.Channel == nil {
		return errors.New("RabbitMQ channel not found when publishing message")
	}

	log.Printf("Publishing message to publisher: %s\n", p.routingKey)

	p.rabbit.Wg.Add(1)
	defer p.rabbit.Wg.Done()

	time.Sleep(10 * time.Second)

	err := p.rabbit.Channel.PublishWithContext(
		ctx,
		"",
		p.routingKey, // queue name
		false,
		false,
		amqp.Publishing{
			Headers:         headers,
			ContentType:     "",
			ContentEncoding: "",
			DeliveryMode:    0,
			Priority:        0,
			CorrelationId:   "",
			ReplyTo:         "",
			Expiration:      "",
			MessageId:       "",
			Timestamp:       time.Now(),
			Type:            "",
			UserId:          "",
			AppId:           "",
			Body:            body,
		},
	)
	// TODO: resend
	if err != nil {
		log.Printf("Couldn't publish message to %s publisher: [%s]\n", p.routingKey, err)
		return err
	}
	return nil
}
