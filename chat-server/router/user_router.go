package router

import (
	v1 "chat-server/api/v1"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

// InitUserRouter 初始化用户相关路由
func (s *UserRouter) InitUserRouter(apiV1 *gin.RouterGroup) {
	// 所有路由都已经通过全局中间件进行JWT认证
	// 这里只需要添加路由即可，无需再添加JWT中间件
	{
		apiV1.POST("/user/register", v1.ApiGroupApp.Register)
		//apiV1.POST("/user/login", v1.ApiGroupApp.Login)
		apiV1.GET("/user/test", v1.ApiGroupApp.Test)
	}
}
