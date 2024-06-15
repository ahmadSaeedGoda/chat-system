package cache

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	DURATION = 30 * time.Second
	DB       = 1
)

var (
	Ctx    = context.Background()
	Client *redis.Client
)

func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   DB,
	})
	pong, err := TestConn(Client, Ctx)
	if err != nil {
		log.Printf("cache is not healthy. returned error: %v", err)
	}
	log.Printf("cache is healthy and responds with '%s'", pong)
}

func Get(key string) (string, error) {
	return Client.Get(Ctx, key).Result()
}

func Set(key string, value interface{}) error {
	return Client.Set(Ctx, key, value, DURATION).Err()
}

func TestConn(client *redis.Client, ctx context.Context) (string, error) {
	// Ping Redis to check connection
	pong, err := client.Ping(ctx).Result()
	return pong, err
}
