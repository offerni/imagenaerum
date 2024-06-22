package img

import "errors"

var (
	ErrInvalidSize error = errors.New("file must not be greater than 5MB")
	ErrNoFiles     error = errors.New("at least one file is required")
	ErrNoSigma     error = errors.New("sigma is required and must be greater than 0")
)
