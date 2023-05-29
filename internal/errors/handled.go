package errors

import (
	"fmt"
	"net/http"
)

type HandledError struct {
	HttpStatus int
	Data       interface{}
}

func (e HandledError) Error() string {
	return fmt.Sprintf("%s", e.Data)
}

func NewValidationError(data interface{}) HandledError {
	return HandledError{
		HttpStatus: http.StatusBadRequest,
		Data:       data,
	}
}

func NewNotFoundError(data interface{}) HandledError {
	return HandledError{
		HttpStatus: http.StatusNotFound,
		Data:       data,
	}
}

func NewUnauthorizedError(data interface{}) HandledError {
	return HandledError{
		HttpStatus: http.StatusUnauthorized,
		Data:       data,
	}
}
