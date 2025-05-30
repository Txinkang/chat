package initialize

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() error {
	redisConfig := AppConfig.Redis // 假设你已经反序列化到 AppConfig.Redis

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Username: redisConfig.Username,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	// 测试连接
	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis 连接失败: %w", err)
	}

	fmt.Println("✅ Redis 连接成功")
	return nil
}
