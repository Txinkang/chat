package service

import (
	"chat-server/global"
	"chat-server/middleware"
	"chat-server/model"
	"chat-server/model/common"
	"github.com/google/uuid"
	"regexp"
	"time"
)

type UserService struct{}

// RegisterUser 注册用户
func (s *UserService) RegisterUser(username, password, email string) (*middleware.TokenPair, error) {
	// 检查用户名是否已存在
	var count int64
	err := global.CHAT_MYSQL.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		global.CHAT_LOG.Error("RegisterUser-->检查用户名，数据库操作错误", "err", err)
		return nil, common.NewServiceError(common.ERROR)
	}
	if count > 0 {
		return nil, common.NewServiceError(common.USERNAME_EXISTS)
	}

	// 判断密码是否合法
	if len(password) <= 0 {
		return nil, common.NewServiceError(common.PASSWORD_INVALID)
	}

	// 判断邮箱是否合法
	if len(email) > 0 && !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
		return nil, common.NewServiceError(common.EMAIL_INVALID)
	}

	// 创建新用户
	userID := uuid.New().String()
	user := model.User{
		ID:        userID,
		Username:  username,
		Password:  password,
		Nickname:  username,
		Email:     email,
		Avatar:    "",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	err = global.CHAT_MYSQL.Create(&user).Error
	if err != nil {
		global.CHAT_LOG.Error("RegisterUser-->创建用户，数据库操作错误", "err", err)
		return nil, common.NewServiceError(common.ERROR)
	}

	// 创建用户成功, 生成token
	tokenPair, err := ServiceGroupApp.TokenService.GenerateTokenPair(userID, username)
	if err != nil {
		global.CHAT_LOG.Error("RegisterUser-->生成token失败", "err", err)
		return nil, common.NewServiceError(common.GENERATE_TOKEN_ERROR)
	}
	// 保存RefreshToken状态
	tokenId := uuid.New().String()
	err = ServiceGroupApp.TokenService.StoreRefreshToken(userID, tokenId, time.Now().Add(time.Hour*time.Duration(global.CHAT_CONFIG.JWT.RefreshTime)), "web")
	if err != nil {
		global.CHAT_LOG.Error("RegisterUser-->保存RefreshToken状态失败", "err", err)
		return nil, common.NewServiceError(common.ERROR)
	}

	return tokenPair, nil
}
