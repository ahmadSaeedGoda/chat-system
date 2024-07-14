.PHONY: default build migrate_up migrate_down dev clean test

default: dev

build:
	docker compose up --build -d

migrate_up:
	migrate -path internal/cassandra/migrations -database cassandra://cassandra:cassandra@localhost:9042/chat up

migrate_down:
	migrate -path internal/cassandra/migrations -database cassandra://cassandra:cassandra@localhost:9042/chat down

dev:
	docker compose up -d

clean:
	docker compose down

test:
	@go test ./...
