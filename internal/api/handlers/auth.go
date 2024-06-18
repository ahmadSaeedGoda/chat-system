package handlers

import (
	"chat-system/internal/api/auth"
	common "chat-system/internal/api/common/constants"
	"chat-system/internal/api/middlewares"
	"chat-system/internal/api/transformers"
	"chat-system/internal/api/validators"
	"chat-system/internal/models"
	"chat-system/internal/services"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	service services.UserService
}

func NewUserHandler(userService services.UserService) *userHandler {
	return &userHandler{
		service: userService,
	}
}

func (s *userHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input models.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		panic(middlewares.NewHTTPError(http.StatusBadRequest, errors.New(common.BAD_REQUEST)))
	}

	if err := validators.ValidateRegisterInput(input); err != nil {
		panic(middlewares.NewHTTPError(http.StatusBadRequest, err))
	}

	exists, err := s.service.UserExists(input.Username)
	if err != nil {
		panic(err)
	}
	if exists {
		panic(middlewares.NewHTTPError(http.StatusConflict, errors.New(common.REGISTER_USER_EXISTS)))
	}

	user, err := s.service.CreateUser(&input)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transformers.TransUserToUserResponse(*user))
}

func (s *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.LoginInput
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		panic(middlewares.NewHTTPError(http.StatusBadRequest, errors.New(common.BAD_REQUEST)))
	}

	user, err := s.service.GetUserByCreds(credentials)

	if err != nil {
		log.Printf("An unexpected error occurred: %v", err)
		panic(middlewares.NewHTTPError(http.StatusUnauthorized, errors.New(common.INVALID_LOGIN)))
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := struct {
		Token string
		Data  transformers.UserResponse
	}{
		Token: token,
		Data:  transformers.TransUserToUserResponse(*user),
	}
	json.NewEncoder(w).Encode(res)
}

/*
TODO:: A list of events/actions that must trigger cache invalidation:
- User Profile Update: When a user updates their profile information (especially username).

- User Deactivation or Suspension: If a user is deactivated or suspended by admin e.g.

- User Logout: When a user logs out, any cached data related to that user should be invalidated to ensure no unauthorized access.

- User Account Deletion: If a user account is deleted.
- Token Refresh Mechanism: If we implement such a mechanism, it's obvious to consider invalidating the respective user cache.
etc...
*/
