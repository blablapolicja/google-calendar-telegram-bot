package repositories

import (
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis"
	"golang.org/x/oauth2"
)

// TokenRepository represents repository for saving and retrieving Google Calendar tokens in Redis
type TokenRepository struct {
	redisClient *redis.Client
	prefix      string
}

// NewTokenRepository creates new TokenRepository
func NewTokenRepository(redisClient *redis.Client) *TokenRepository {
	prefix := "ID_"

	return &TokenRepository{redisClient, prefix}
}

// Save serializes and saves token in Redis
func (r *TokenRepository) Save(ID int64, token *oauth2.Token) error {
	serialized, err := json.Marshal(token)

	if err != nil {
		return err
	}

	return r.redisClient.Set(r.prefix+strconv.FormatInt(ID, 10), serialized, 0).Err()
}

// Get gets token from Redis and deserializes it
func (r *TokenRepository) Get(ID int64) (*oauth2.Token, error) {
	serializedToken, err := r.redisClient.Get(r.prefix + strconv.FormatInt(ID, 10)).Result()

	if err != nil {
		return nil, err
	}

	var token *oauth2.Token

	if err := json.Unmarshal([]byte(serializedToken), &token); err != nil {
		return nil, err
	}

	return token, nil
}
