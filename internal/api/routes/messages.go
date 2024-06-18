package routes

import (
	"chat-system/internal/api/middlewares"

	"github.com/gorilla/mux"
)

func getMsgsRoutes(apiRouter *mux.Router) *mux.Router {
	msgRouter := apiRouter.PathPrefix("/messages").Subrouter().StrictSlash(true)

	// Apply Auth middleware
	msgRouter.Use(middlewares.IsAuth)

	msgRouter.HandleFunc("/send", appConfig.GetMsgHandler().SendMessage).Methods("POST")
	msgRouter.HandleFunc("/", appConfig.GetMsgHandler().GetMessages).Methods("GET")

	return apiRouter
}
