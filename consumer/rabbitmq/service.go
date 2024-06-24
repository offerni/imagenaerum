package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Service struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

type NewServiceOpts struct {
	Url string
}

func NewService(opts NewServiceOpts) (*Service, error) {
	conn, err := amqp.Dial(opts.Url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	log.Println("RabbitMQ Channel opened")

	return &Service{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (s *Service) NewChannel() (*amqp.Channel, error) {
	return s.Conn.Channel()
}
