package consumer

import (
	"context"
	"errors"
	"github.com/vadimpk/image-resizing-server-go/pkg/queue/rabbitmq"
	"log"
)

type Consumer struct {
	rabbit     *rabbitmq.RabbitMQ
	routingKey string
	cfg        *rabbitmq.QueueConfig
}

func NewConsumer(rabbit *rabbitmq.RabbitMQ) (*Consumer, error) {

	cfg, err := rabbitmq.Init("configs/queues")
	if err != nil {
		log.Println("couldn't init rabbitmq/consumer config")
		return nil, err
	}

	// TODO: try to create new
	if rabbit.Channel == nil {
		return nil, errors.New("RabbitMQ channel not found when creating producer")
	}

	q, err := rabbit.Channel.QueueDeclare(
		cfg.Name,
		cfg.Durable,
		cfg.DeleteUnused,
		cfg.Exclusive,
		cfg.NoWait,
		nil,
	)
	if err != nil {
		log.Printf("couldn't create RabbitMQ queue: [%s]\n", err)
		return nil, err
	}

	err = rabbit.Channel.Qos(cfg.PrefetchCount, 0, false)
	if err != nil {
		log.Printf("couldn't set Qos: [%s]\n", err)
		return nil, err
	}

	c := &Consumer{
		rabbit:     rabbit,
		routingKey: q.Name,
		cfg:        cfg,
	}

	return c, nil
}

func (c *Consumer) Close(ctx context.Context) chan struct{} {
	return c.rabbit.Close(ctx)
}

// Consume consumes messages from channel and applies func f to each message
// if context gets cancelled, it tries to process last message and then stops channel
func (c *Consumer) Consume(ctx context.Context, f func(body []byte, headers map[string]interface{})) error {

	// TODO: create new channel (reconnect)
	if c.rabbit.Channel == nil {
		return errors.New("RabbitMQ channel not found when consuming messages")
	}

	c.rabbit.Wg.Add(1)
	defer c.rabbit.Wg.Done()

	msgs, err := c.rabbit.Channel.Consume(
		c.routingKey,
		"",
		c.cfg.AutoAck,
		c.cfg.Exclusive,
		false,
		c.cfg.NoWait,
		nil,
	)

	if err != nil {
		return err
	}
	var allCanceled bool
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				return errors.New("error receiving message from channel ")
			}
			log.Println("Consumed message")

			c.rabbit.Wg.Add(1)
			f(msg.Body, msg.Headers)
			if err := msg.Ack(false); err != nil {
				log.Printf("error acking message: %s\n", err)
			}
			c.rabbit.Wg.Done()

		case <-ctx.Done():
			// if context is cancelled, try to process last message, then stop
			if allCanceled {
				return nil
			}
			err = c.rabbit.Channel.Cancel(c.routingKey, false) // stop receiving
			allCanceled = true
			continue
		}
	}
}
