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
)

const CACHE_KEY_SUFFIX = "-messages"

type MsgHandler interface {
	SendMessage(w http.ResponseWriter, r *http.Request)
	GetMessages(w http.ResponseWriter, r *http.Request)
}

type msgHandler struct {
	service     services.MessageService
	userService services.UserService
}

func NewMsgHandler(msgService services.MessageService, userService services.UserService) *msgHandler {
	return &msgHandler{
		service:     msgService,
		userService: userService,
	}
}

func (mh *msgHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var input models.SendMessageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		panic(middlewares.NewHTTPError(http.StatusBadRequest, errors.New(common.BAD_REQUEST)))
	}

	if err := validators.ValidateSendMessageInput(input); err != nil {
		panic(middlewares.NewHTTPError(http.StatusBadRequest, errors.New(common.BAD_REQUEST)))
	}

	exists, err := mh.userService.UserExists(input.Recipient)
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
	if err := mh.service.CreateMessage(msg); err != nil {
		panic(err)
	}

	senderCacheKey := userClaims.Username + CACHE_KEY_SUFFIX

	if err := mh.service.UpdateCachedMsgsForUser(senderCacheKey, *msg); err != nil {
		panic(err)
	}

	recipientCacheKey := msg.Recipient + CACHE_KEY_SUFFIX

	// Let's do the same for recipient
	if err := mh.service.UpdateCachedMsgsForUser(recipientCacheKey, *msg); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

// GetMessages retrieves all messages for the authenticated user
func (mh *msgHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	var (
		messages       []models.Message
		err            error
		page, pageSize int
	)

	userClaims := middlewares.GetUserFromContext(r.Context())
	username := userClaims.Username
	page, pageSize = utils.GetPaginationParams(r)

	cacheKey := username + CACHE_KEY_SUFFIX
	messages, err = mh.service.GetFromCache(cacheKey)
	if err != nil {
		panic(err)
	}

	if messages == nil {
		messages, err = mh.service.GetMessages(username)
		if err != nil {
			panic(err)
		}

		err = mh.service.SetMessagesToCache(cacheKey, messages)
		if err != nil {
			panic(err)
		}
	}

	res := paginateMessages(page, pageSize, messages)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// TODO:: When message update, delete cache must be invalidated

// TODO:: Can be more generic to paginate any entity
// paginateMessages paginates the messages
func paginateMessages(page, pageSize int, messages []models.Message) map[string]interface{} {
	itemsCount := len(messages)

	startIndex := (page - 1) * pageSize

	if startIndex > itemsCount {
		startIndex = itemsCount
	}

	endIndex := page * pageSize
	if endIndex > itemsCount {
		endIndex = itemsCount
	}

	paginatedMessages := messages[startIndex:endIndex]

	res := map[string]interface{}{
		"messages": paginatedMessages,
		"pagination": map[string]interface{}{
			"currentPage":   page,
			"pageSize":      pageSize,
			"totalMessages": itemsCount,
			"totalPages":    (itemsCount + pageSize - 1) / pageSize,
		},
	}

	return res
}
