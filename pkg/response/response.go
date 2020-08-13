package response

import (
	"encoding/json"
	"net/http"
)

//ErrorMessage is a message
type ErrorMessage struct {
	Message string `json:"message"`
}

//HTTPError work for error
func HTTPError(w http.ResponseWriter, r *http.Request, statusCode int, message string) error {
	msg := ErrorMessage{
		Message: message,
	}

	j, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(j)
	return nil
}
