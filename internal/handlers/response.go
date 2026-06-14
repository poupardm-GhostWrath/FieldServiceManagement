// internal/handlers/response.go
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithError(w http.ResponseWriter, code int, message string, err error) {
	type response struct {
		Message string `json:"message"`
	}

	msg := message

	if err != nil {
		msg = fmt.Sprintf("%s: %s", message, err.Error())
	}

	RespondWithJSON(w, code, response{
		Message: msg,
	})
}
