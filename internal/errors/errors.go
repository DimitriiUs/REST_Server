package errors

import "errors"

var (
	ErrNotFound           = errors.New("not found task")
	ErrInvalidID          = errors.New("invalid id")
	ErrInvalidDescription = errors.New("invalid description")
	ErrInvalidDueDate     = errors.New("invalid due date")
)
