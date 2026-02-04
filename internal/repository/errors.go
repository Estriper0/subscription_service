package repository

import "errors"

const (
	PgCodeConstrainError = "23514"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrIncorrectTime = errors.New("the end date must be later than the start date")
)
