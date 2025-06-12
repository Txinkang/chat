package service

import (
	"chat-server/global"
	"chat-server/middleware"
	"chat-server/model"
	"chat-server/model/common"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

// RegisterUser 注册用户
func (s *UserService) RegisterUser(userAccount, password, email string) (*middleware.TokenPair, error) {
	// 检查用户名是否已存在
	var count int64
	err := global.CHAT_MYSQL.Model(&model.User{}).Where("user_account = ?", userAccount).Count(&count).Error
	if err != nil {
		global.CHAT_LOG.Error("RegisterUser-->检查用户账号，数据库操作错误", "err", err)
		return nil, common.NewServiceError(common.ERROR)
	}
	if count > 0 {
		return nil, common.NewServiceError(common.USER_ACCOUNT_EXISTS)
	}

	// 判断密码是否合法
	if len(password) <= 0 {
		return nil, common.NewServiceError(common.PASSWORD_INVALID)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		global.CHAT_LOG.Error("RegisterUser-->加密密码出错", "err", err)
		return nil, common.NewServiceError(common.ERROR)
	}

	// 判断邮箱是否合法
	if len(email) > 0 && !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
		return nil, common.NewServiceError(common.EMAIL_INVALID)
	}

	// 创建新用户
	userID := uuid.New().String()
	user := model.User{
		ID:          userID,
		UserAccount: userAccount,
		Password:    string(hashedPassword),
		Nickname:    userAccount,
		Email:       email,
		Avatar:      "",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	err = global.CHAT_MYSQL.Create(&user).Error
	if err != nil {
		global.CHAT_LOG.Error("RegisterUser-->创建用户，数据库操作错误", "err", err)
		return nil, common.NewServiceError(common.ERROR)
	}

	// 创建用户成功, 生成token
	tokenPair, err := ServiceGroupApp.TokenService.GenerateTokenPair(userID, userAccount)
	if err != nil {
		global.CHAT_LOG.Error("RegisterUser-->生成token失败", "err", err)
		return nil, common.NewServiceError(common.GENERATE_TOKEN_ERROR)
	}
	// 在redis保存RefreshToken状态
	tokenId := uuid.New().String()
	err = ServiceGroupApp.TokenService.StoreRefreshToken(userID, tokenId, time.Now().Add(time.Hour*24*time.Duration(global.CHAT_CONFIG.JWT.RefreshTime)), "web")
	if err != nil {
		global.CHAT_LOG.Error("RegisterUser-->保存RefreshToken状态失败", "err", err)
		return nil, common.NewServiceError(common.ERROR)
	}

	return tokenPair, nil
}

func (s *UserService) LoginAccount(userAccount string, password string, platform string) (*middleware.TokenPair, error) {
	// 验证userAccount
	var queryUser model.User
	err := global.CHAT_MYSQL.Where("user_account = ?", userAccount).First(&queryUser).Error
	if err != nil {
		global.CHAT_LOG.Error("LoginAccount-->检查用户账号，数据库操作错误", "err", err)
		return nil, common.NewServiceError(common.USER_ACCOUNT_NOT_FOUND)
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(queryUser.Password), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, common.NewServiceError(common.PASSWORD_INVALID)
	}

	// 检查redis是否已存在该登录平台的token，只允许单平台登录
	tokenKey := fmt.Sprintf("user_tokens:%s", queryUser.ID)
	tokenIds, err := global.CHAT_REDIS.SMembers(context.Background(), tokenKey).Result()
	if err != nil {
		global.CHAT_LOG.Error("LoginAccount-->检查redis是否已存在该登录平台的token失败", "err", err)
		return nil, common.NewServiceError(common.ERROR)
	}
	if len(tokenIds) > 0 {
		for _, tokenId := range tokenIds {
			// 通过tokenId获取refreshToken
			refreshToken := fmt.Sprintf("refresh_token:%s:%s", queryUser.ID, tokenId)
			refreshTokenData, err := global.CHAT_REDIS.Get(context.Background(), refreshToken).Result()
			if err != nil {
				global.CHAT_LOG.Error("LoginAccount-->获取refreshToken失败", "err", err)
				return nil, common.NewServiceError(common.ERROR)
			}
			// 解析值，获取登录平台信息
			var tokenData map[string]interface{}
			if err = json.Unmarshal([]byte(refreshTokenData), &tokenData); err != nil {
				global.CHAT_LOG.Error("LoginAccount-->解析refreshTokenData失败", "err", err)
				return nil, common.NewServiceError(common.ERROR)

			}
			// 平台相同则撤销旧令牌
			getPlatform, ok := tokenData["platform"].(string)
			if !ok {
				global.CHAT_LOG.Error("LoginAccount-->获取platform失败", "err", err)
				return nil, common.NewServiceError(common.ERROR)
			}
			if getPlatform == platform {
				err := ServiceGroupApp.TokenService.RevokeToken(queryUser.ID, tokenId)
				if err != nil {
					return nil, err
				}
			}

		}
	}
	// 通过验证，下发token
	tokenPair, err := ServiceGroupApp.TokenService.GenerateTokenPair(queryUser.ID, queryUser.UserAccount)
	if err != nil {
		global.CHAT_LOG.Error("LoginAccount-->生成token失败", "err", err)
		return nil, common.NewServiceError(common.GENERATE_TOKEN_ERROR)
	}

	// 在redis保存RefreshToken状态
	tokenId := uuid.New().String()
	err = ServiceGroupApp.TokenService.StoreRefreshToken(queryUser.ID, tokenId, time.Now().Add(time.Hour*24*time.Duration(global.CHAT_CONFIG.JWT.RefreshTime)), "web")
	if err != nil {
		global.CHAT_LOG.Error("LoginAccount-->保存RefreshToken状态失败", "err", err)
		return nil, common.NewServiceError(common.ERROR)
	}

	return tokenPair, nil
}
