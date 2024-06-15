package validators

import (
	"chat-system/internal/models"
)

func ValidateSendMessageInput(input models.SendMessageInput) error {
	return validate.Struct(input)
}
