package core

import (
	"chat-server/global"
	"context"
	"log/slog"
	"time"
)

func CloseResource() {

	if global.CHAT_MYSQL != nil {
		slog.Info("关闭 MySQL 连接...")
		if sqlDB, err := global.CHAT_MYSQL.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				slog.Error("关闭 MySQL 连接失败:", "err", err)
			} else {
				slog.Info("MySQL 连接已关闭")
			}
		} else {
			slog.Error("获取 GORM 底层 *sql.DB 失败，无法关闭 MySQL:", "err", err)
		}
	}

	if global.CHAT_REDIS != nil {
		slog.Info("关闭 Redis 连接...")
		if err := global.CHAT_REDIS.Close(); err != nil {
			slog.Error("关闭 Redis 连接失败:", "err", err)
		} else {
			slog.Info("Redis 连接已关闭")
		}
	}

	if global.CHAT_MONGO != nil {
		slog.Info("关闭 MongoDB 连接...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := global.CHAT_MONGO.Disconnect(ctx); err != nil {
			slog.Error("关闭 MongoDB 连接失败:", "err", err)
		} else {
			slog.Info("MongoDB 连接已关闭")
		}
	}
}
