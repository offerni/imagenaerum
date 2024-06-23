package img

import "errors"

var (
	ErrInvalidSize       error = errors.New("file must not be greater than 5MB")
	ErrNoFiles           error = errors.New("at least one file is required")
	ErrNoRabbitMQService error = errors.New("rabbitMQ Service is required")
	ErrNoParam           error = errors.New("param is required")
)
