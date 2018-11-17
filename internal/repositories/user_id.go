package repositories

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// UserIDRepository represents repository for saving and retrieving userIDs which need to be authorized
type UserIDRepository struct {
	redisClient *redis.Client
	prefix      string
}

// NewUserIDRepository creates new UserIDRepository
func NewUserIDRepository(redisClient *redis.Client) *UserIDRepository {
	prefix := "STATE_"

	return &UserIDRepository{redisClient, prefix}
}

// Save userID
func (r *UserIDRepository) Save(state string, userID int64) error {
	return r.redisClient.Set(r.prefix+state, strconv.FormatInt(userID, 10), time.Minute).Err()
}

// Get userID
func (r *UserIDRepository) Get(state string) (int64, error) {
	value, err := r.redisClient.Get(r.prefix + state).Result()

	if err != nil {
		return 0, err
	}

	userID, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

// Delete userID
func (r *UserIDRepository) Delete(state string) error {
	return r.redisClient.Del(r.prefix + state).Err()
}
