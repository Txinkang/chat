package service

import (
	"chat-server/global"
	"chat-server/model"
	"errors"
	"github.com/google/uuid"
)

type UserService struct{}

// RegisterUser 注册用户
func (s *UserService) RegisterUser(username, password, email string) error {
	// 检查用户名是否已存在
	var count int64
	err := global.CHAT_MYSQL.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}
	userID := uuid.New().String()
	user := model.User{
		ID:       userID,
		Username: username,
		Password: password,
		Email:    email,
	}
	err = global.CHAT_MYSQL.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}
