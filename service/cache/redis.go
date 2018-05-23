package cache

import "github.com/go-redis/redis"

var RedisClient *redis.Client

func InitializeRedisClient(address string, password string) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})
}
