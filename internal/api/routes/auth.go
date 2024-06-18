package routes

import (
	"github.com/gorilla/mux"
)

func getAuthRoutes(apiRouter *mux.Router) *mux.Router {
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", appConfig.GetUserHandler().Register).Methods("POST")
	authRouter.HandleFunc("/login", appConfig.GetUserHandler().Login).Methods("POST")
	return authRouter
}
