package router

import (
	v1 "chat-server/api/v1"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

// InitUserRouter 初始化用户相关路由
func (s *UserRouter) InitUserRouter(apiV1 *gin.RouterGroup) {
	{
		apiV1.POST("/user/register", v1.ApiGroupApp.Register)
		//apiV1.POST("/user/login", v1.ApiGroupApp.Login)
		apiV1.GET("/user/test", v1.ApiGroupApp.Test)
	}
}
