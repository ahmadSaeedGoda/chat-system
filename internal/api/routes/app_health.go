package routes

import (
	"chat-system/internal/api/handlers"

	"github.com/gorilla/mux"
)

func getAppRoutes(apiRouter *mux.Router) *mux.Router {
	apiRouter.HandleFunc("/health", handlers.HowAmI).Methods("GET")

	return apiRouter
}
