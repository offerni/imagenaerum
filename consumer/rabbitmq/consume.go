package rabbitmq

import (
	"context"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumeOpts struct {
	Ch       *amqp.Channel
	Out      chan<- amqp.Delivery
	QName    string
	Consumer string
}

func (svc Service) Consume(opts ConsumeOpts) error {
	msgs, err := opts.Ch.ConsumeWithContext(
		context.Background(),
		opts.QName,
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
