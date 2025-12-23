package domain

import "errors"

var (
	ErrNotFound     = errors.New("record not found")
	ErrUnauthorized = errors.New("unauthorized action")
	ErrForbidden    = errors.New("forbidden action")
	ErrBadRequest   = errors.New("bad request")
)
