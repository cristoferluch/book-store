package http_errors

import "net/http"

const (
	ErrBookNotFound = "book not found"
	Unexpected      = "something went wrong"
)

type HandlerError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewUnexpectedError(message string) HandlerError {
	return HandlerError{message, http.StatusInternalServerError}
}

func NewNotFoundError(message string) HandlerError {
	return HandlerError{message, http.StatusNotFound}
}

func NewForbiddenError(message string) HandlerError {
	return HandlerError{message, http.StatusForbidden}
}

func NewBadRequestError(message string) HandlerError {
	return HandlerError{message, http.StatusBadRequest}
}

func (e HandlerError) Error() string {
	return e.Message
}
