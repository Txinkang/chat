package router

import (
	"github.com/gin-gonic/gin"
)

type ChatRouter struct{}

// InitChatRouter 初始化聊天相关路由
func (s *ChatRouter) InitChatRouter(apiV1 *gin.RouterGroup) {
	// 聊天相关路由 - 需要认证

}
