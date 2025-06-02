package v1

import (
	"chat-server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatApi struct{}

// 发送消息请求结构
type SendMessageRequest struct {
	Content string `json:"content" binding:"required"`
	ToUser  string `json:"to_user" binding:"required"`
}

// SendMessage 发送聊天消息
func (i *ChatApi) SendMessage(c *gin.Context) {
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 从上下文获取当前用户ID (假设中间件已经设置)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未认证",
		})
		return
	}

	// 调用service层处理发送消息逻辑
	err := service.SendMessage(userID.(string), req.ToUser, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "发送消息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "发送成功",
	})
}

// GetChatHistory 获取聊天历史
func (i *ChatApi) GetChatHistory(c *gin.Context) {
	// 获取查询参数
	toUser := c.Query("to_user")
	if toUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "缺少to_user参数",
		})
		return
	}

	// 从上下文获取当前用户ID (假设中间件已经设置)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未认证",
		})
		return
	}

	// 调用service层获取聊天历史
	messages, err := service.GetChatHistory(userID.(string), toUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取聊天历史失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"msg":      "获取成功",
		"messages": messages,
	})
}
