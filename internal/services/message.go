package services

import (
	"chat-system/internal/api/cache"
	"chat-system/internal/db"
	"chat-system/internal/models"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gocql/gocql"
)

type MessageService interface {
	CreateMessage(message *models.Message) error
	GetMessages(username string) ([]models.Message, error)
	GetFromCache(cacheKey string) ([]models.Message, error)
	SetMessagesToCache(cacheKey string, messages []models.Message) error
}

type messageService struct{}

func NewMessageService() *messageService {
	return &messageService{}
}

func (s *messageService) CreateMessage(message *models.Message) error {
	message.ID = gocql.TimeUUID()
	message.Timestamp = time.Now()

	query := `INSERT INTO chat.messages
		(user, timestamp, id, sender, recipient, content)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	batch := db.Session.NewBatch(gocql.LoggedBatch)

	batch.Query(
		query,
		message.Sender,
		message.Timestamp,
		message.ID,
		message.Sender,
		message.Recipient,
		message.Content,
	)

	batch.Query(
		query,
		message.Recipient,
		message.Timestamp,
		message.ID,
		message.Sender,
		message.Recipient,
		message.Content,
	)

	return db.Session.ExecuteBatch(batch)
}

func (s *messageService) GetMessages(username string) ([]models.Message, error) {
	var messages []models.Message
	iter := db.Session.Query(
		`SELECT id, sender, recipient, timestamp, content
		FROM chat.messages
		WHERE user = ?
		ORDER BY timestamp DESC`,
		username,
	).Iter()
	var message models.Message
	for iter.Scan(&message.ID, &message.Sender, &message.Recipient, &message.Timestamp, &message.Content) {
		messages = append(messages, message)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *messageService) GetFromCache(cacheKey string) ([]models.Message, error) {
	jsonData, err := cache.Get(cacheKey)
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var messages []models.Message
	err = json.Unmarshal([]byte(jsonData), &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *messageService) SetMessagesToCache(cacheKey string, messages []models.Message) error {
	jsonData, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	err = cache.Set(cacheKey, jsonData)
	if err != nil {
		return err
	}

	return nil
}
