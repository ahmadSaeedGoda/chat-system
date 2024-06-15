package routes

import (
	"chat-system/internal/api/handlers"

	"github.com/gorilla/mux"
)

func getAuthRoutes(apiRouter *mux.Router) *mux.Router {
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", handlers.Register).Methods("POST")
	authRouter.HandleFunc("/login", handlers.Login).Methods("POST")
	return authRouter
}
