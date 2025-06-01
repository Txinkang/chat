package initialize

import (
	"chat-server/service" // 导入 service 包
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync" // 导入 sync 包
)

// Initialize 函数负责所有应用程序的初始化
// 它接受一个 context 和 WaitGroup，用于统一管理 Goroutine 生命周期
func Initialize(appCtx context.Context, appCancel context.CancelFunc, wg *sync.WaitGroup) error { // 接受 appCtx 和 wg

	// 1. 初始化配置文件 (检查错误)
	if err := InitConfig(); err != nil {
		return fmt.Errorf("初始化配置文件失败: %w", err)
	}

	// 2. 初始化日志 (检查错误)
	if err := InitLogger(); err != nil {
		return fmt.Errorf("初始化日志工具失败: %w", err)
	}

	// 3. 初始化数据库 (检查错误)
	if err := InitMySQL(); err != nil {
		return fmt.Errorf("初始化 MySQL 失败: %w", err)
	}
	if err := InitRedis(); err != nil {
		return fmt.Errorf("初始化 Redis 失败: %w", err)
	}
	if err := InitMongo(); err != nil {
		return fmt.Errorf("初始化 MongoDB 失败: %w", err)
	}
	if err := InitElasticSearch(); err != nil {
		return fmt.Errorf("初始化 Elasticsearch 失败: %w", err)
	}

	// 4. 启动数据同步服务 (或其他后台服务)
	wg.Add(1) // 增加 Goroutine 计数
	go func() {
		defer wg.Done() // Goroutine 结束时递减计数
		defer func() {
			if r := recover(); r != nil {
				slog.Error("数据同步 Goroutine 发生 panic", "panic_value", r, "stack_trace", string(debug.Stack()))
			}
		}()

		syncErr := service.StartMongoToEsSync(appCtx, "messages", "messages")
		if syncErr != nil && syncErr != context.Canceled {
			slog.Error("MongoDB 到 Elasticsearch 同步服务终止，发生非取消错误", "err", syncErr)
			appCancel()
		} else if syncErr == context.Canceled {
			slog.Info("MongoDB 到 Elasticsearch 同步服务已因 Context 取消而停止。")
		}
	}()

	slog.Info("所有应用程序组件初始化完成。")
	return nil
}
