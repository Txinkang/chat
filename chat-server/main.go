package main

import (
	"chat-server/core"
	"chat-server/initialize"
	"log/slog"
	"os"

	"context"
	"sync"
)

func main() {
	//设置全局context，用来统一关闭goroutine
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()
	var wg sync.WaitGroup

	err := initialize.Initialize(appCtx, appCancel, &wg)
	if err != nil {
		slog.Error("应用程序初始化失败:", "err", err)
		os.Exit(1)
	}
	core.RunServe(appCtx, appCancel, &wg)
	//最后关闭各种资源
	defer core.CloseResource()
}
