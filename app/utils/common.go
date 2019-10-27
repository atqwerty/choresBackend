package utils

import (
	"encoding/json"
	"net/http"
)

// RespondJSON ...
func RespondJSON(w http.ResponseWriter, status int, data interface{}) error {
	response, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
	return nil
}

// ServerError ...
func ServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Server encountered an error." + err.Error()))
}

// BadRequest ...
func BadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}
