package service

import (
	"time"
)

type ChatService struct{}

// 生成消息ID
func (s *ChatService) generateMessageID() string {
	return "msg_" + time.Now().Format("20060102150405") + "_" + s.randomString(8)
}

// 生成随机字符串
func (s *ChatService) randomString(length int) string {
	// 简化实现，实际应使用更安全的随机数生成
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[time.Now().UnixNano()%int64(len(chars))]
		time.Sleep(1 * time.Nanosecond) // 确保每次获取不同的时间戳
	}
	return string(result)
}
