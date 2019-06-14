package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type status string

const (
	ok          = status("OK")
	badRequest  = status("BadRequest")
	serverError = status("ServerError")
)

type response struct {
	Status  status      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func handleServerError(w http.ResponseWriter, err error) {
	jsonResponse, _ := json.Marshal(response{
		Status:  serverError,
		Message: err.Error(),
	})

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonResponse)
}

// SecretPage to test auth
func SecretPage(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(ContextUserKey("username")).(string)
	w.Write([]byte(fmt.Sprintf("Hello %s", username)))
}
