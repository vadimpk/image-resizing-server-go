package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type Producer struct {
	amqpChan   *amqp.Channel
	amqpConn   *amqp.Connection
	routingKey string
}

func NewProducer(conn *amqp.Connection, queueConfig QueueConfig) (*Producer, error) {
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Couldn't open RabbitMQ channel: [%s]\n", err)
		return nil, err
	}

	q, err := ch.QueueDeclare(
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
		amqpChan:   ch,
		amqpConn:   conn,
		routingKey: q.Name,
	}

	return p, nil
}

func (p *Producer) Publish(ctx context.Context, pub amqp.Publishing) {

	log.Printf("Publishing message to publisher: %s\n", p.routingKey)

	err := p.amqpChan.PublishWithContext(
		ctx,
		"",           // we don't use exchange
		p.routingKey, // publisher name
		false,
		false,
		pub,
	)
	// TODO: resend
	if err != nil {
		log.Printf("Couldn't publish message to %s publisher: [%s]\n", p.routingKey, err)
	}
}

func (p *Producer) PublishImage(ctx context.Context, body []byte, headers map[string]interface{}) {

	pub := amqp.Publishing{
		Headers:         headers,
		ContentType:     "",
		ContentEncoding: "",
		DeliveryMode:    amqp.Persistent,
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
	}

	p.Publish(ctx, pub)
}
