package service

import (
	"errors"
	"time"
)

// Message 消息模型
type Message struct {
	ID        string    `json:"id"`
	FromUser  string    `json:"from_user"`
	ToUser    string    `json:"to_user"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// 模拟消息存储
var messageStore = make([]Message, 0)

// SendMessage 发送消息
func SendMessage(fromUser, toUser, content string) error {
	// 检查发送者和接收者是否存在
	if _, exists := userStore[fromUser]; !exists {
		return errors.New("发送者不存在")
	}
	if _, exists := userStore[toUser]; !exists {
		return errors.New("接收者不存在")
	}

	// 创建消息
	message := Message{
		ID:        generateMessageID(),
		FromUser:  fromUser,
		ToUser:    toUser,
		Content:   content,
		CreatedAt: time.Now(),
	}

	// 存储消息
	messageStore = append(messageStore, message)

	return nil
}

// GetChatHistory 获取聊天历史
func GetChatHistory(userID, otherUserID string) ([]Message, error) {
	// 检查两个用户是否存在
	if _, exists := userStore[userID]; !exists {
		return nil, errors.New("用户不存在")
	}
	if _, exists := userStore[otherUserID]; !exists {
		return nil, errors.New("对方用户不存在")
	}

	// 筛选出两个用户之间的消息
	var chatHistory []Message
	for _, msg := range messageStore {
		if (msg.FromUser == userID && msg.ToUser == otherUserID) ||
			(msg.FromUser == otherUserID && msg.ToUser == userID) {
			chatHistory = append(chatHistory, msg)
		}
	}

	return chatHistory, nil
}

// 生成消息ID
func generateMessageID() string {
	return "msg_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

// 生成随机字符串
func randomString(length int) string {
	// 简化实现，实际应使用更安全的随机数生成
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[time.Now().UnixNano()%int64(len(chars))]
		time.Sleep(1 * time.Nanosecond) // 确保每次获取不同的时间戳
	}
	return string(result)
}
