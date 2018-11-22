package redisstorage

import (
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/config"
	"github.com/go-redis/redis"
)

// NewRedisClient creates new Redis client
func NewRedisClient(config config.RedisConf) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.Database,
	})

	_, err := redisClient.Ping().Result()

	return redisClient, err
}
