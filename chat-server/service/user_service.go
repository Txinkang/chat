package service

import (
	"errors"
	"time"
)

// 模拟用户数据存储
var userStore = make(map[string]User)

// User 用户模型
type User struct {
	Username string
	Password string
	Email    string
}

// RegisterUser 注册用户
func RegisterUser(username, password, email string) error {
	// 检查用户名是否已存在
	if _, exists := userStore[username]; exists {
		return errors.New("用户名已存在")
	}

	// 存储用户信息 (实际应用中应该对密码进行哈希处理)
	userStore[username] = User{
		Username: username,
		Password: password,
		Email:    email,
	}

	return nil
}

// LoginUser 用户登录
func LoginUser(username, password string) (string, error) {
	// 检查用户是否存在
	user, exists := userStore[username]
	if !exists {
		return "", errors.New("用户不存在")
	}

	// 验证密码 (实际应用中应该比较哈希值)
	if user.Password != password {
		return "", errors.New("密码错误")
	}

	// 生成JWT令牌 (这里简化处理，实际应该使用JWT库)
	token := "jwt_token_" + username + "_" + time.Now().Format("20060102150405")

	return token, nil
}
