package img

import "github.com/offerni/imagenaerum/consumer/rabbitmq"

type Service struct {
	RabbitMQService *rabbitmq.Service
}

type NewServiceOpts struct {
	RabbitMQSvc *rabbitmq.Service
}

func NewService(opts NewServiceOpts) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return &Service{
		RabbitMQService: opts.RabbitMQSvc,
	}, nil
}

func (opts NewServiceOpts) Validate() error {
	if opts.RabbitMQSvc == nil {
		return ErrNoRabbitMQService
	}

	return nil
}
