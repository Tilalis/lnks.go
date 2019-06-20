package main

import (
	"flag"
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
	var configFile string

	const (
		defaultFilename = "config.json"
		flagUsage       = "configuration file"
	)

	flag.StringVar(&configFile, "config", defaultFilename, flagUsage)
	flag.StringVar(&configFile, "c", defaultFilename, flagUsage)
	flag.Parse()

	cfg, err := config.LoadConfig(configFile)

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
		"/api/create",
		auth.WeakMiddleware(handlers.CreateAlias),
	).Methods("POST")

	router.HandleFunc(
		"/api/delete",
		auth.StrongMiddleware(handlers.DeleteAlias),
	).Methods("POST")

	router.HandleFunc(
		"/api/all",
		auth.StrongMiddleware(handlers.GetAliases),
	).Methods("GET")

	router.HandleFunc(
		"/api/{alias}",
		handlers.GetAlias,
	).Methods("GET")

	router.HandleFunc("/{alias}", handlers.ResolveAlias).Methods("GET")

	log.Fatal(http.ListenAndServe(cfg.Address, router))
}
