package errs

import "errors"

var (
	// Common errors
	ErrInternal          = errors.New("internal server error")
	ErrNotImplementedYet = errors.New("not implemented yet")
	ErrNilNotAllowed     = errors.New("nil values not allowed")
	ErrTimeout           = errors.New("timeout")
	ErrInvalidRequest    = errors.New("invalid request")
	ErrNotFound          = errors.New("not found")
	ErrBadId             = errors.New("bad id")
	ErrBadUserId         = errors.New("bad user id")

	// Specific errors
	ErrInvalidInputField = errors.New("invalid input")
)
