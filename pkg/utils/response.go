package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type APIResponse struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func SendJSON(w http.ResponseWriter, payload interface{}, code int) {

	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(payload)
	if err != nil {
		slog.Error("failed to marshal response",
			slog.String("error", err.Error()),
			slog.String("payload", payload.(string)),
		)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"something went wrong"}`))
		return
	}

	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		slog.Error("failed to write response",
			slog.String("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"something went wrong"}`))
	}
}
