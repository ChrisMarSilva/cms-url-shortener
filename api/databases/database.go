package databases

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func CreateClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_REDIS_ADDR"),
		Password: os.Getenv("DB_REDIS_PASS"),
		DB:       dbNo,
	})

	return rdb
}
