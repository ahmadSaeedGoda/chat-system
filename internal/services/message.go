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

type MessageService struct{}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) CreateMessage(message *models.Message) error {
	message.ID = gocql.TimeUUID()
	message.Timestamp = time.Now()

	queryRecipient := `INSERT INTO chat.messages_by_recipient (recipient, id, sender, timestamp, content) VALUES (?, ?, ?, ?, ?)`
	querySender := `INSERT INTO chat.messages_by_sender (sender, id, recipient, timestamp, content) VALUES (?, ?, ?, ?, ?)`

	batch := db.Session.NewBatch(gocql.LoggedBatch)

	batch.Query(
		queryRecipient,
		message.Recipient,
		message.ID,
		message.Sender,
		message.Timestamp,
		message.Content,
	)

	batch.Query(
		querySender,
		message.Sender,
		message.ID,
		message.Recipient,
		message.Timestamp,
		message.Content,
	)

	return db.Session.ExecuteBatch(batch)
}

func (s *MessageService) GetMessages(username string) ([]models.Message, error) {
	var messages []models.Message
	iter := db.Session.Query(`SELECT id, sender, recipient, timestamp, content FROM messages WHERE recipient = ?`, username).Iter()
	var message models.Message
	for iter.Scan(&message.ID, &message.Sender, &message.Recipient, &message.Timestamp, &message.Content) {
		messages = append(messages, message)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *MessageService) GetFromCache(cacheKey string) ([]models.Message, error) {
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

func (s *MessageService) SetMessagesToCache(cacheKey string, messages []models.Message) error {
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

func (s *MessageService) GetFromDB(username string, page, pageSize int) ([]models.Message, error) {
	var (
		messages                   []models.Message
		id                         gocql.UUID
		sender, recipient, content string
		timestamp                  time.Time
	)

	seenMessages := make(map[gocql.UUID]bool)

	// Fetch messages where the user is the recipient
	queryRecipient := `SELECT recipient, id, sender, timestamp, content FROM chat.messages_by_recipient WHERE recipient = ?`
	iterRecipient := db.Session.Query(queryRecipient, username).Iter()

	for iterRecipient.Scan(&recipient, &id, &sender, &timestamp, &content) {
		if !seenMessages[id] {
			messages = append(
				messages,
				models.Message{
					ID:        id,
					Sender:    sender,
					Recipient: recipient,
					Timestamp: timestamp,
					Content:   content,
				},
			)
			seenMessages[id] = true
		}
	}
	if err := iterRecipient.Close(); err != nil {
		return nil, err
	}

	// Fetch messages where the user is the sender
	querySender := `SELECT sender, id, recipient, timestamp, content FROM chat.messages_by_sender WHERE sender = ?`
	iterSender := db.Session.Query(querySender, username).Iter()

	for iterSender.Scan(&sender, &id, &recipient, &timestamp, &content) {
		if !seenMessages[id] {
			messages = append(
				messages, models.Message{
					ID:        id,
					Sender:    sender,
					Recipient: recipient,
					Timestamp: timestamp,
					Content:   content,
				},
			)
			seenMessages[id] = true
		}
	}
	if err := iterSender.Close(); err != nil {
		return nil, err
	}

	return messages, nil
}
