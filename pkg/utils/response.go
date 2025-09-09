package utils

import (
	"book-store/pkg/utils/http_errors"
	"errors"
	"github.com/bytedance/sonic"
	"log/slog"
	"net/http"
)

type apiResponse struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func SendJSON(w http.ResponseWriter, payload interface{}, code int) {

	w.Header().Set("Content-Type", "application/json")

	var p apiResponse
	p.Data = payload

	response, err := sonic.Marshal(p)
	if err != nil {
		slog.Error("failed to marshal response",
			slog.String("error", err.Error()),
			slog.Any("payload", payload),
		)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"something went wrong"}`))
		return
	}

	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func SendError(w http.ResponseWriter, err error) {

	var handlerError http_errors.HandlerError

	var code int
	var message string

	if errors.As(err, &handlerError) {
		code = handlerError.Code
		message = handlerError.Message
	} else {
		code = http.StatusInternalServerError
		message = "something went wrong"
	}

	var p apiResponse
	p.Error = message

	w.Header().Set("Content-Type", "application/json")

	response, err := sonic.Marshal(p)
	if err != nil {
		slog.Error("failed to marshal response",
			slog.String("error", err.Error()),
			slog.Any("payload", p),
		)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"something went wrong"}`))
		return
	}

	w.WriteHeader(code)
	_, _ = w.Write(response)
}
