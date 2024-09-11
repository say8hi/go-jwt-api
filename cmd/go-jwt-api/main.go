package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/say8hi/go-jwt-api/internal/database"
	"github.com/say8hi/go-jwt-api/internal/handlers"
	"github.com/say8hi/go-jwt-api/internal/middlewares"
)

func main() {
	database.Init()
	database.CreateTables()
	defer database.CloseConnection()

	r := mux.NewRouter()
	r.Use(middlewares.LoggingMiddleware)

	// authRouter := r.NewRoute().Subrouter()
	// authRouter.Use(middlewares.AuthMiddleware)

	// Unauthorized endpoints
	// Users
	r.HandleFunc("/users/create", handlers.CreateUserHandler).Methods("POST")
	r.HandleFunc("/users/get_tokens", handlers.CreateTokensPair).Methods("GET")
	r.HandleFunc("/users/refresh", handlers.RefreshTokens).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
