package appconfig

import (
	"chat-system/internal/api/handlers"
	dbmanager "chat-system/internal/db_manager"
	"chat-system/internal/services"
)

type AppConfig interface {
	GetUserHandler() handlers.AuthHandler
	GetMsgHandler() handlers.MsgHandler
}

type appConfig struct {
}

func NewAppConfig() *appConfig {
	return &appConfig{}
}

func (a *appConfig) GetUserHandler() handlers.AuthHandler {
	return handlers.NewUserHandler(
		services.NewUserService(
			dbmanager.CassandraSession,
			dbmanager.CASSANDRA_KEYSPACE,
			dbmanager.USERS_TABLE,
		),
	)
}

func (a *appConfig) GetMsgHandler() handlers.MsgHandler {
	return handlers.NewMsgHandler(
		services.NewMessageService(
			dbmanager.CassandraSession,
			dbmanager.CASSANDRA_KEYSPACE,
			dbmanager.MSGS_TABLE,
		),
		services.NewUserService(
			dbmanager.CassandraSession,
			dbmanager.CASSANDRA_KEYSPACE,
			dbmanager.USERS_TABLE,
		),
	)
}
