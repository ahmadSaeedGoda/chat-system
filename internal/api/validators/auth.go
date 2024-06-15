package validators

import (
	"chat-system/internal/models"

	validator "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateRegisterInput(input models.RegisterInput) error {
	return validate.Struct(input)
}

func ValidateLoginInput(input models.LoginInput) error {
	return validate.Struct(input)
}
