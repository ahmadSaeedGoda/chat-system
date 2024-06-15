package routes

import (
	"chat-system/internal/api/handlers"
	"chat-system/internal/api/middlewares"

	"github.com/gorilla/mux"
)

func getMsgsRoutes(apiRouter *mux.Router) *mux.Router {
	msgRouter := apiRouter.PathPrefix("/messages").Subrouter().StrictSlash(true)

	// Apply Auth middleware
	msgRouter.Use(middlewares.IsAuth)

	msgRouter.HandleFunc("/send", handlers.SendMessage).Methods("POST")
	msgRouter.HandleFunc("/", handlers.GetMessages).Methods("GET")

	return apiRouter
}
