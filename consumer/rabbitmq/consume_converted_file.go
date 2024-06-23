package rabbitmq

import (
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (svc *Service) ConsumeConvertedFile(queueName string) error {
	_, err := svc.Channel.QueueDeclare(
		queueName,
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
		Ch:       svc.Channel,
		Out:      msgsChannel,
		QName:    queueName,
		Consumer: "img-consumer",
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
