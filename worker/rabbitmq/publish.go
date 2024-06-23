package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type PublishOpts struct {
	Body         []byte
	Ch           *amqp.Channel
	ExchangeName string
	QueueName    string
	RoutingKey   string
}

func (svc Service) Publish(opts PublishOpts) error {
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

	err = opts.Ch.Publish(
		opts.ExchangeName,
		opts.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        opts.Body,
		})
	if err != nil {
		return err
	}

	return nil
}
