package errs

import (
	"errors"
	"net/http"
)

// Constant Errors
var (
	ErrNotEnoughBalance = StatusError{
		Err:    errors.New("not enough balance"),
		Status: http.StatusBadRequest,
	}

	ErrNoUser = StatusError{
		Err:    errors.New("the user not found"),
		Status: http.StatusBadRequest,
	}
	ErrNoItem = StatusError{
		Err:    errors.New("the item do not exist"),
		Status: http.StatusBadRequest,
	}
	ErrUnAuth = StatusError{
		Err:    errors.New("the user is not authorized"),
		Status: http.StatusUnauthorized,
	}
	//ErrInvalidReq       = errors.New("invalid request body")
)

/*
200-ok
unauthorized
bad request
internal
*/

type StatusError struct {
	Err    error
	Status int
}

func WrapError(err error, status int) StatusError {
	return StatusError{
		Err:    err,
		Status: status,
	}
}

func (s StatusError) Error() string {
	return s.Err.Error()
}
