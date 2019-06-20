package handlers

import (
	"encoding/json"
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
