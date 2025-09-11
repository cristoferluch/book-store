package http_error

import "net/http"

const (
	ErrBookNotFound       = "book not found"
	ErrUnexpected         = "something went wrong"
	ErrInvalidId          = "invalid id"
	ErrInvalidRequestBody = "invalid request body"
	ErrInvalidQueryParams = "invalid query params"

	ErrBookAuthorIsRequired = "book author is required"
	ErrBookIsbnIsRequired   = "book isbn is required"
	ErrBookTitleIsRequired  = "book title is required"
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
