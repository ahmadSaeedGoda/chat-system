package services

import (
	"chat-system/internal/api/cache"
	"chat-system/internal/cassandra"
	"chat-system/internal/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gocql/gocql"
)

type MessageService interface {
	CreateMessage(message *models.Message) error
	GetMessages(username string) ([]models.Message, error)
	GetFromCache(cacheKey string) ([]models.Message, error)
	SetMessagesToCache(cacheKey string, messages []models.Message) error
	UpdateCachedMsgsForUser(cacheKey string, msg models.Message) error
}

type messageService struct {
	db         *gocql.Session
	dbKeyspace string
	tableName  string
}

func NewMessageService(db *gocql.Session, keyspace, tableName string) *messageService {
	return &messageService{
		db:         db,
		dbKeyspace: keyspace,
		tableName:  tableName,
	}
}

func (s *messageService) CreateMessage(message *models.Message) error {
	message.ID = gocql.TimeUUID()
	message.Timestamp = time.Now()

	query := fmt.Sprintf(
		`INSERT INTO %s.%s
		(user, timestamp, id, sender, recipient, content)
		VALUES (?, ?, ?, ?, ?, ?)`,
		s.dbKeyspace,
		s.tableName,
	)

	batch := s.db.NewBatch(gocql.LoggedBatch)

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

	return cassandra.Session.ExecuteBatch(batch)
}

func (s *messageService) GetMessages(username string) ([]models.Message, error) {
	var messages []models.Message
	query := fmt.Sprintf(
		`SELECT id, sender, recipient, timestamp, content
		FROM %s.%s
		WHERE user = ?
		ORDER BY timestamp DESC`,
		s.dbKeyspace,
		s.tableName,
	)
	iter := cassandra.Session.Query(
		query,
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

func (s *messageService) UpdateCachedMsgsForUser(cacheKey string, msg models.Message) error {
	// Append the message to the Redis list
	// Set the updated messages list in Redis
	messages, err := s.GetFromCache(cacheKey)
	if err != nil {
		return err
	}

	// if messages are not cached, we cannot update the cache now
	// return early, fail fast and wait for messages to be available in cache to update them!
	if messages == nil {
		return nil
	}

	// Prepend to cached since they are already sorted by DB
	messages = append([]models.Message{msg}, messages...)
	err = s.SetMessagesToCache(cacheKey, messages)
	if err != nil {
		return err
	}

	return nil
}
