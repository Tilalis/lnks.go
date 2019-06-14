package handlers

import (
	"context"
	"encoding/json"
	"lnks/config"
	"lnks/models"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// ContextUserKey context user key
type ContextUserKey string

func getToken(username, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})
	return token.SignedString([]byte(secretKey))
}

func verifyToken(tokenString, secretKey string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	username, ok := token.Claims.(jwt.MapClaims)["username"]

	if !ok {
		return nil, ErrWrongToken
	}

	user, _ := models.GetUser(username.(string))
	return user, err
}

// Auth struct
type Auth struct {
	Config *config.Config
}

// Authenticate authentication handler; requires config.Inject
func (a *Auth) Authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var authenticationRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&authenticationRequest)

	if err != nil {
		handleServerError(w, err)
		return
	}

	user, err := models.AuthenticateUser(
		authenticationRequest.Username,
		authenticationRequest.Password,
	)

	if err != nil {
		var message string

		if err == models.ErrWrongUserPassword {
			message = "Wrong user password"
		} else if user == nil {
			message = "There is no such user"
		}

		jsonResponse, _ := json.Marshal(response{
			Status:  badRequest,
			Message: message,
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)

		return
	}

	token, err := getToken(user.Username, a.Config.SecretKey)

	if err != nil {
		handleServerError(w, err)
		return
	}

	jsonResponse, _ := json.Marshal(response{
		Status: ok,
		Data:   token,
	})

	w.Header().Set("Authorization", "Bearer "+token)
	w.Write(jsonResponse)
}

// Middleware auth middleware
func (a *Auth) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
		user, err := verifyToken(tokenString, a.Config.SecretKey)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(""))
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserKey("username"), user.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
