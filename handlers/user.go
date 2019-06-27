package handlers

import (
	"encoding/json"
	"lnks/models"
	"net/http"
)

// RegisterUser handler
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request authRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)

	if err != nil {
		handleServerError(w, err)
		return
	}

	user, err := models.NewUser(request.Username, request.Password)

	if err != nil {
		jsonResponse, _ := json.Marshal(response{
			Status:  badRequest,
			Message: "Empty credentials are not allowed",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}

	err = user.Save()

	if err != nil {
		handleServerError(w, err)
		return
	}

	jsonResponse, _ := json.Marshal(response{
		Status:  ok,
		Message: "Successfully registered user",
	})

	w.Write(jsonResponse)
}
