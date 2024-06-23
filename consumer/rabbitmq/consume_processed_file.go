package rabbitmq

import (
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumeProcessedFileOpts struct {
	ExchangeName string
	QueueName    string
	RoutingKey   string
}

func (svc *Service) ConsumeProcessedFile(opts ConsumeProcessedFileOpts) error {
	_, err := svc.Channel.QueueDeclare(
		opts.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgsChannel := make(chan amqp.Delivery)
	if err := svc.Consume(ConsumeOpts{
		Ch:           svc.Channel,
		Consumer:     "img-consumer",
		ExchangeName: opts.ExchangeName,
		Out:          msgsChannel,
		QueueName:    opts.QueueName,
		RoutingKey:   opts.RoutingKey,
	}); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for msg := range msgsChannel {
		wg.Add(1)
		go func(m amqp.Delivery) {
			defer wg.Done()
			log.Printf("Received a message: %s", m.Body)
			m.Ack(false)
		}(msg)
	}

	wg.Wait()
	return nil
}
