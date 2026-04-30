package go_redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewClient(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis ping %v", err)
	}
	return rdb
}
