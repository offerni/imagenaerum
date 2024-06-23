package rabbitmq

import (
	"context"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumeOpts struct {
	Ch           *amqp.Channel
	Consumer     string
	ExchangeName string
	Out          chan<- amqp.Delivery
	QueueName    string
	RoutingKey   string
}

func (svc Service) Consume(opts ConsumeOpts) error {
	msgs, err := opts.Ch.ConsumeWithContext(
		context.Background(),
		opts.QueueName,
		opts.Consumer,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = svc.Channel.ExchangeDeclare(
		opts.ExchangeName,
		"direct",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = svc.Channel.QueueBind(
		opts.QueueName,
		opts.RoutingKey,
		opts.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	go func() {
		defer close(opts.Out)
		for msg := range msgs {
			wg.Add(1)
			go func(m amqp.Delivery) {
				defer wg.Done()
				opts.Out <- m
			}(msg)
		}
	}()
	wg.Wait()

	return nil
}
