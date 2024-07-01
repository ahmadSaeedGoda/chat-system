package middlewares

import (
	common "chat-system/internal/api/common/constants"
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse is a struct for sending JSON error responses
type ErrorResponse struct {
	Error string `json:"error"`
}

// HTTPError is a custom error type that includes an HTTP status code
type HTTPError struct {
	StatusCode int
	Err        error
}

// Error is a method to get the error message from HTTPError object
func (e *HTTPError) Error() string {
	return e.Err.Error()
}

// NewHTTPError creates a new HTTPError object
func NewHTTPError(statusCode int, err error) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Err:        err,
	}
}

// ErrorHandlerMiddleware is a middleware that wraps all routes to catch errors
// and send an appropriate JSON response
func HandleErrors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				switch e := err.(type) {
				case *HTTPError:
					log.Printf("HTTP error: %v", e.Err)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(e.StatusCode)
					json.NewEncoder(w).Encode(ErrorResponse{Error: e.Error()})
				default:
					log.Printf("An unexpected error occurred: %v", err)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(ErrorResponse{Error: common.INTERNAL_SERVER_ERROR})
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
