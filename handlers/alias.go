package handlers

import (
	"encoding/json"
	"lnks/models"
	"net/http"

	"github.com/gorilla/mux"
)

// ResolveAlias view for resolving aliases
func ResolveAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["alias"]

	alias, err := models.GetAlias(name)

	if alias == nil {
		http.Redirect(w, r, "/", 301)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, alias.URL, 301)
}

// CreateAlias view for creating aliases
func CreateAlias(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	alias := &models.Alias{}
	decoder := json.NewDecoder(r.Body)

	// Decode
	err := decoder.Decode(alias)

	if err != nil {
		handleServerError(w, err)
		return
	}

	// Validate
	err = alias.Validate()

	if err != nil {
		var resp *response

		switch err {
		case models.ErrWrongAlias:
			resp = &response{
				Status:  badRequest,
				Message: "Alias must contain only alphanumeric characters.",
			}
		case models.ErrWrongURL:
			resp = &response{
				Status:  badRequest,
				Message: "Malformed URL",
			}
		}

		jsonResponse, _ := json.Marshal(resp)

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)

		return
	}

	// Save

	// Get user from auth middleware
	if user, ok := r.Context().Value(ContextKey("user")).(*models.User); ok {
		alias.SetUser(user)
	}

	err = alias.Save()

	if err != nil {
		if err == models.ErrNoConnection {
			handleServerError(w, err)
			return
		}

		jsonResponse, _ := json.Marshal(response{
			Status:  badRequest,
			Message: "Alias already exists.",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}

	jsonResponse, _ := json.Marshal(response{
		Status:  ok,
		Message: "Alias created",
		Data:    alias,
	})

	w.Write(jsonResponse)
}

// DeleteAlias view for deleting alias
func DeleteAlias(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	alias := &models.Alias{}
	decoder := json.NewDecoder(r.Body)

	// Decode
	err := decoder.Decode(alias)

	if err != nil {
		handleServerError(w, err)
		return
	}

	err = alias.Delete()

	if err != nil {
		handleServerError(w, err)
	}

	jsonResponse, _ := json.Marshal(response{
		Status:  ok,
		Message: "Alias deleted.",
	})

	w.Write(jsonResponse)
}

// GetAlias view for getting single alias (for testing)
func GetAlias(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	name := vars["alias"]

	alias, _ := models.GetAlias(name)

	if alias == nil {
		jsonResponse, _ := json.Marshal(response{
			Status:  badRequest,
			Message: "There is no such alias",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}

	jsonResponse, _ := json.Marshal(response{
		Status: ok,
		Data:   alias,
	})

	w.Write(jsonResponse)
}

// GetAliases view for getting all aliases for specific user
func GetAliases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(ContextKey("user")).(*models.User)
	aliases, err := models.GetAliases(user)

	if err != nil {
		handleServerError(w, err)
		return
	}

	jsonResponse, _ := json.Marshal(response{
		Status:  ok,
		Message: "",
		Data:    aliases,
	})

	w.Write(jsonResponse)
}
