package core

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// RunServe 负责启动应用程序并处理优雅关闭逻辑
// 它会阻塞当前 Goroutine，直到收到操作系统信号或 context 被取消
func RunServe(appCtx context.Context, appCancel context.CancelFunc, wg *sync.WaitGroup) {
	// 阻塞主 Goroutine，等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	slog.Info("应用程序启动完成，等待退出信号...")

	// 使用 select 监听信号或 appCtx 的取消
	select {
	case sig := <-sigChan:
		slog.Info("收到退出信号，开始优雅关闭...", "signal", sig.String())
	case <-appCtx.Done():
		slog.Info("应用程序 Context 已被取消，开始优雅关闭...", "context_err", appCtx.Err())
	}

	// 调用 appCancel() 来通知所有子 Goroutine 停止工作
	appCancel()

	// 等待所有 Goroutine 完成它们的清理工作
	waitTimeout := 10 * time.Second // 适当延长等待时间，例如 10 秒
	done := make(chan struct{})

	// 启动一个 Goroutine 来等待 WaitGroup 归零
	go func() {
		wg.Wait()   // 阻塞直到所有 Goroutine 都调用了 Done()
		close(done) // 通知 wg.Wait() 已完成
	}()

	// 使用 select 监听 Goroutine 是否完成或等待超时
	select {
	case <-done:
		slog.Info("所有 Goroutine 已安全退出。")
	case <-time.After(waitTimeout):
		slog.Warn("等待 Goroutine 退出超时，强制关闭。")
	}

	slog.Info("应用程序已关闭。")
}
