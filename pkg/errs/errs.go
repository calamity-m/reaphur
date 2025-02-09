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
	ErrBadId             = errors.New("id not formed correctly")
	ErrBadUserId         = errors.New("user id not formed correctly")

	// Specific errors
	ErrInvalidInputField = errors.New("invalid input")
)
