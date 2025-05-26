package authRedis

import (
	"github.com/go-redis/redis"
	"time"
)

// Kurinov
type Casher struct {
	*redis.Client
	RefreshTTL time.Duration
}

// NewRedisClient returns new redis client
func NewRedisClient(address string, password string, dbID int, refreshTTL time.Duration) *Casher {

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       dbID,
	})
	return &Casher{
		Client:     client,
		RefreshTTL: refreshTTL}
}
