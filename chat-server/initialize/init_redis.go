package initialize

import (
	"chat-server/global"
	"context"
	"log/slog"
	"time"

	"github.com/go-redis/redis/v8"
)

func InitRedis() error {
	slog.Info("初始化redis")

	redisConfig := global.CHAT_CONFIG.Redis // 假设你已经反序列化到 AppConfig.Redis
	global.CHAT_REDIS = redis.NewClient(&redis.Options{
		Addr:         redisConfig.Address,
		Username:     redisConfig.Username,
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
		PoolSize:     redisConfig.PoolSize,                                 // 连接池大小
		MinIdleConns: redisConfig.MinIdleConns,                             // 最小空闲连接数
		PoolTimeout:  time.Duration(redisConfig.PoolTimeout) * time.Second, // 获取连接的超时时间
		IdleTimeout:  time.Duration(redisConfig.IdleTimeout) * time.Second, // 空闲连接超时时间
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := global.CHAT_REDIS.Ping(ctx).Err(); err != nil {
		slog.Error("Redis 连接失败: ", "err", err)
		closeErr := global.CHAT_REDIS.Close()
		if closeErr != nil {
			slog.Error("Redis 连接失败后，关闭 Redis 失败: ", "err", closeErr)
		}
		return err
	}
	slog.Info("Redis 连接成功")
	return nil
}
