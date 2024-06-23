package img

import "github.com/offerni/imagenaerum/worker/rabbitmq"

type Service struct {
	RabbitmqService *rabbitmq.Service
}

type NewServiceOpts struct {
	RabbitmqSvc *rabbitmq.Service
}

func NewService(opts NewServiceOpts) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return &Service{
		RabbitmqService: opts.RabbitmqSvc,
	}, nil
}

func (opts NewServiceOpts) Validate() error {
	if opts.RabbitmqSvc == nil {
		return ErrNoRabbitMQSvc
	}
	return nil
}
