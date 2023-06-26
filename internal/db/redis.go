package db

import (
	"context"
	"fmt"
	"time"

	"github.com/SideSwapIN/Analystic/internal/config"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedisDB() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         config.GetConfig().Redis.Host,
		Password:     config.GetConfig().Redis.Password,
		DB:           config.GetConfig().Redis.DB,
		PoolSize:     10,               // 设置连接池大小
		MinIdleConns: 5,                // 设置最小空闲连接数
		IdleTimeout:  30 * time.Second, // 设置连接最长空闲时间
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("failed to connect Redis: %v", err)
	}
	return nil
}
