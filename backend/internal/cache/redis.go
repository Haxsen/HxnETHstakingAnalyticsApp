package cache

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

// InitRedis initializes the Redis client connection
func InitRedis() error {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}

	Client = redis.NewClient(opt)

	// Test the connection
	ctx := context.Background()
	_, err = Client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	log.Println("Redis connection established")
	return nil
}

// CloseRedis closes the Redis connection
func CloseRedis() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// Set sets a key-value pair in Redis with expiration
func Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value from Redis by key
func Get(ctx context.Context, key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

// Delete removes a key from Redis
func Delete(ctx context.Context, key string) error {
	return Client.Del(ctx, key).Err()
}

// Exists checks if a key exists in Redis
func Exists(ctx context.Context, key string) (bool, error) {
	result, err := Client.Exists(ctx, key).Result()
	return result > 0, err
}
