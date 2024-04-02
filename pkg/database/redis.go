package database

import (
	"github.com/go-redis/redis/v8"
	"os"
)

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RDB_ADDR"),
		Password: os.Getenv("RDB_PASSWORD"),
		DB:       0,
	})
	return rdb
}
