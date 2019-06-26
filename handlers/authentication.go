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

// ContextKey context user key
type ContextKey string

// Auth struct
type Auth struct {
	cfg *config.Config
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewAuth auth constructor
func NewAuth(cfg *config.Config) *Auth {
	if cfg == nil {
		panic("lnks: auth config cannot be nil")
	}

	return &Auth{
		cfg: cfg,
	}
}

func (auth *Auth) getToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})
	return token.SignedString([]byte(auth.cfg.SecretKey))
}

func (auth *Auth) verifyToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(auth.cfg.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	username, ok := token.Claims.(jwt.MapClaims)["username"]

	if !ok {
		return nil, ErrWrongToken
	}

	return models.GetUser(username.(string))
}

// Authenticate authentication handler
func (auth *Auth) Authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request authRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)

	if err != nil {
		handleServerError(w, err)
		return
	}

	user, err := models.AuthenticateUser(
		request.Username,
		request.Password,
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

	token, err := auth.getToken(user.Username)

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

func (auth *Auth) middleware(next http.HandlerFunc, strict bool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context

		authorizationHeader := r.Header.Get("Authorization")

		if !strict && authorizationHeader == "" {
			ctx = context.WithValue(r.Context(), ContextKey("user"), nil)
		} else {
			tokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)
			user, err := auth.verifyToken(tokenString)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(""))
				return
			}

			ctx = context.WithValue(r.Context(), ContextKey("user"), user)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// StrictMiddleware auth middleware which does not allow empty Authorization
func (auth *Auth) StrictMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return auth.middleware(next, true)
}

// Middleware auth middleware which allows empty Authorization
func (auth *Auth) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return auth.middleware(next, false)
}
