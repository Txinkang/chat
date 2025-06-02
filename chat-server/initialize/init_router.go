package initialize

import (
	"chat-server/global"
	"chat-server/middleware"
	"chat-server/router"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化所有路由
func InitRouter() {
	Router := gin.New()
	Router.Use(gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		Router.Use(gin.Logger())
	}
	Router.Use(middleware.Cors())
	apiV1 := Router.Group("/api/v1")

	// 初始化各个模块的路由
	router.RouterGroupApp.InitChatRouter(apiV1)
	router.RouterGroupApp.InitUserRouter(apiV1)

	global.CHAT_ROUTERS = Router

}
