package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func ConnectRedis(ctx context.Context, host, port, password string, db int) error {
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := Redis.Ping(pingCtx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Redis connected successfully")
	return nil
}

func SetRedisKey(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return Redis.Set(ctx, key, value, expiration).Err()
}

func GetRedisKey(ctx context.Context, key string) (string, error) {
	return Redis.Get(ctx, key).Result()
}

func DelRedisKey(ctx context.Context, key string) error {
	return Redis.Del(ctx, key).Err()
}

func CloseRedisConn() {
	if Redis != nil {
		Redis.Close()
	}
}
