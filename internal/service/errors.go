package service

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrInternal      = errors.New("internal error")
	ErrIncorrectTime = errors.New("the end date must be later than the start date")
)
