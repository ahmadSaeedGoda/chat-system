package transformers

import (
	"chat-system/internal/models"

	"github.com/gocql/gocql"
)

// UserResponse struct with only the fields to be exposed
type UserResponse struct {
	ID       gocql.UUID `json:"id"`
	Username string     `json:"username"`
}

// Function to map User to UserResponse
func TransUserToUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}
}
