package errs

import "errors"

var (
	// Common errors
	ErrInternal          = errors.New("internal server error")
	ErrNotImplementedYet = errors.New("not implemented yet")
	ErrNilNotAllowed     = errors.New("nil values not allowed")
	ErrTimeout           = errors.New("timeout")
	ErrNotFound          = errors.New("not found")
	ErrBadRequest        = errors.New("bad request")
	ErrBadId             = errors.New("bad id")
	ErrBadUserId         = errors.New("bad user id")

	// Specific errors
	ErrInvalidInputField = errors.New("invalid input")
)
