package rest

import "errors"

var (
	ErrInvalidForm error = errors.New("unable to parse multipart form")
	ErrEmptyForm   error = errors.New("form cannot be empty")
)
