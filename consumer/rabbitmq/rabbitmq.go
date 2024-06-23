package rabbitmq

import "log"

func Start() *Service {
	rabbitMQSvc, err := NewService(NewServiceOpts{
		Url: Url,
	})
	if err != nil {
		log.Fatalf("%s: %s", "Failed to initialize RabbitMQ service", err)
	}

	return rabbitMQSvc
}

func (s *Service) Close() {
	if s.Channel != nil {
		s.Channel.Close()
	}
	if s.Conn != nil {
		s.Conn.Close()
	}
}
