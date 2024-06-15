package models

import (
	"time"

	"github.com/gocql/gocql"
)

type User struct {
	ID       gocql.UUID `json:"id"`
	Username string     `json:"username"`
	Password string     `json:"-"`
}

type RegisterInput struct {
	Username string `json:"username" validate:"required,min=1,max=16"`
	Password string `json:"password" validate:"required,min=6,max=12"`
}

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Message struct {
	ID        gocql.UUID `json:"id"`
	Sender    string     `json:"sender"`
	Recipient string     `json:"recipient"`
	Timestamp time.Time  `json:"timestamp"`
	Content   string     `json:"content" validate:"required,min=1,max=1000"`
}

type SendMessageInput struct {
	Content   string `json:"content" validate:"required,min=1,max=1000"`
	Recipient string `json:"recipient" validate:"required,min=1,max=16"`
}
