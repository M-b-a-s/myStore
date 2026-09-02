package json

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Write(w http.ResponseWriter, status int, data any, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}

	_ = json.NewEncoder(w).Encode(response)
}

func Read(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
