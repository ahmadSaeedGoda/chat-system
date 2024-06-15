package services

import (
	"chat-system/internal/db"
	"chat-system/internal/models"
	"errors"

	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetUserByCreds(credentials models.LoginInput) (*models.User, error) {
	var existingUser models.User
	err := db.Session.Query(`SELECT id, username, password FROM users WHERE username = ? LIMIT 1`, credentials.Username).
		Scan(&existingUser.ID, &existingUser.Username, &existingUser.Password)

	if err != nil {
		return nil, errors.New("invalid login")
	}

	if !s.checkPasswordHash(credentials.Password, existingUser.Password) {
		return nil, errors.New("invalid login")
	}

	return &existingUser, err
}

func (s *UserService) UserExists(username string) (bool, error) {
	var existingUserId gocql.UUID
	err := db.Session.Query(`SELECT id FROM users WHERE username = ? LIMIT 1`, username).Scan(&existingUserId)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, gocql.ErrNotFound) {
		return false, nil
	}

	return false, err
}

func (s *UserService) CreateUser(userInput *models.RegisterInput) (*models.User, error) {
	hashedPassword, err := s.hashPassword(userInput.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:       gocql.TimeUUID(),
		Username: userInput.Username,
		Password: hashedPassword,
	}

	err = db.Session.Query(
		`INSERT INTO users (id, username, password) VALUES (?, ?, ?)`,
		user.ID, user.Username, user.Password,
	).Exec()

	return user, err
}

func (s *UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *UserService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
