package routes

import (
	"chat-system/internal/api/middlewares"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	r := mux.NewRouter()

	// Apply the error handler middleware
	r.Use(middlewares.HandleErrors)

	// API routes
	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	getAppRoutes(apiRouter)
	getAuthRoutes(apiRouter)
	getMsgsRoutes(apiRouter)

	return r
}
