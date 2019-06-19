package main

import (
	"lnks/config"
	"lnks/handlers"
	"lnks/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}

func main() {
	cfg, err := config.LoadConfig("config.json")

	if err != nil {
		log.Fatal(err)
		return
	}

	models.Prepare(cfg.Database.Driver, cfg.Database.ConnectionString)
	defer models.Close()

	auth := handlers.NewAuth(cfg)
	router := mux.NewRouter()

	router.HandleFunc("/", index).Methods("GET")

	router.HandleFunc(
		"/api/login",
		auth.Authenticate,
	).Methods("POST")

	router.HandleFunc(
		"/api/secretpage",
		auth.Middleware(handlers.SecretPage),
	).Methods("GET")

	router.HandleFunc(
		"/api/create",
		handlers.CreateAlias,
	).Methods("POST")

	router.HandleFunc(
		"/api/all",
		handlers.GetAliases,
	).Methods("GET")

	router.HandleFunc(
		"/api/{alias}",
		handlers.GetAlias,
	).Methods("GET")

	router.HandleFunc("/{alias}", handlers.ResolveAlias).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
