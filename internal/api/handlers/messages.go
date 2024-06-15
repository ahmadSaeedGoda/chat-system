package handlers

import (
	common "chat-system/internal/api/common/constants"
	"chat-system/internal/api/common/utils"
	"chat-system/internal/api/middlewares"
	"chat-system/internal/api/validators"
	"chat-system/internal/models"
	"chat-system/internal/services"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
)

const CACHE_KEY_SUFFIX = "-messages"

var msgService = services.NewMessageService()

func SendMessage(w http.ResponseWriter, r *http.Request) {
	var input models.SendMessageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		panic(middlewares.NewHTTPError(http.StatusBadRequest, errors.New(common.BAD_REQUEST)))
	}

	if err := validators.ValidateSendMessageInput(input); err != nil {
		panic(middlewares.NewHTTPError(http.StatusBadRequest, errors.New(common.BAD_REQUEST)))
	}

	exists, err := userService.UserExists(input.Recipient)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic(middlewares.NewHTTPError(http.StatusBadRequest, errors.New(common.SEND_MESSAGE_NO_RECIPIENT)))
	}

	userClaims := middlewares.GetUserFromContext(r.Context())

	msg := &models.Message{
		Sender:    userClaims.Username,
		Recipient: input.Recipient,
		Content:   input.Content,
	}
	if err := msgService.CreateMessage(msg); err != nil {
		panic(err)
	}

	cacheKey := userClaims.Username + CACHE_KEY_SUFFIX

	// Append the message to the Redis list
	messages, err := msgService.GetFromCache(cacheKey)
	if err != nil {
		panic(err)
	}

	messages = append(messages, *msg)

	// Set the updated messages list in Redis
	err = msgService.SetMessagesToCache(cacheKey, messages)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	var (
		messages       []models.Message
		err            error
		page, pageSize int
	)

	userClaims := middlewares.GetUserFromContext(r.Context())
	username := userClaims.Username
	page, pageSize = utils.GetPaginationParams(r)

	cacheKey := username + CACHE_KEY_SUFFIX
	messages, err = msgService.GetFromCache(cacheKey)
	if err != nil {
		panic(err)
	}

	if messages == nil {
		messages, err = msgService.GetFromDB(username, page, pageSize)
		if err != nil {
			panic(err)
		}

		err = msgService.SetMessagesToCache(cacheKey, messages)
		if err != nil {
			panic(err)
		}
	}

	// Sort messages by timestamp in descending order
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Timestamp.After(messages[j].Timestamp)
	})

	res := paginateMessages(page, pageSize, messages)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// TODO:: When message update, delete cache must be invalidated

// TODO:: Can be more generic to paginate any entity
// paginateMessages paginates the messages
func paginateMessages(page, pageSize int, messages []models.Message) map[string]interface{} {
	startIndex := (page - 1) * pageSize
	endIndex := page * pageSize
	if endIndex > len(messages) {
		endIndex = len(messages)
	}

	paginatedMessages := messages[startIndex:endIndex]

	res := map[string]interface{}{
		"messages": paginatedMessages,
		"pagination": map[string]interface{}{
			"currentPage":   page,
			"pageSize":      pageSize,
			"totalMessages": len(messages),
			"totalPages":    (len(messages) + pageSize - 1) / pageSize,
		},
	}

	return res
}
