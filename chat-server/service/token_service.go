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
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenService struct{}

// GenerateTokenPair 生成JWT令牌
func (s *TokenService) GenerateTokenPair(userID string, userAccount string) (*middleware.TokenPair, error) {
	// 生成唯一的令牌ID
	tokenID := uuid.New().String()

	// 创建访问令牌
	accessClaims := middleware.AccessToken{
		UserID:      userID,
		UserAccount: userAccount,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(global.CHAT_CONFIG.JWT.AccessTime))), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                                     // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                                                     // 生效时间
			Issuer:    global.CHAT_CONFIG.JWT.Issuer,                                                                      // 签发人
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(global.CHAT_CONFIG.JWT.Secret))
	if err != nil {
		global.CHAT_LOG.Error("GenerateTokenPair-->签名生成accessToken失败", "err", err)
		return nil, common.NewServiceError(common.ERROR)
	}

	// 创建刷新令牌
	refreshClaims := middleware.RefreshToken{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(global.CHAT_CONFIG.JWT.RefreshTime))), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                                         // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                                                         // 生效时间
			Issuer:    global.CHAT_CONFIG.JWT.Issuer,                                                                          // 签发人
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(global.CHAT_CONFIG.JWT.Secret))
	if err != nil {
		global.CHAT_LOG.Error("GenerateTokenPair-->签名生成refreshToken失败", "err", err)
		return nil, common.NewServiceError(common.ERROR)
	}

	return &middleware.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    global.CHAT_CONFIG.JWT.AccessTime,
	}, nil
}

// ParseAccessToken 解析访问令牌
func (s *TokenService) ParseAccessToken(tokenString string) (*middleware.AccessToken, error) {
	claims := &middleware.AccessToken{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.CHAT_CONFIG.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshAccessToken 刷新访问令牌
func (s *TokenService) RefreshAccessToken(refreshTokenString string) (string, error) {
	refreshClaims := &middleware.RefreshToken{}
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.CHAT_CONFIG.JWT.Secret), nil
	})

	if err != nil || !refreshToken.Valid {
		return "", common.NewServiceError(common.REFRESH_TOKEN_INVALID)
	}

	// 检查令牌是否被撤销
	isTokenRevoked, err := s.isTokenRevoked(refreshClaims.UserID, refreshClaims.TokenID)
	if err != nil {
		return "", common.NewServiceError(common.ERROR)
	}
	if isTokenRevoked {
		return "", common.NewServiceError(common.REFRESH_TOKEN_REVOKED)
	}

	// 获取用户信息
	user, err := s.getUserByID(refreshClaims.UserID)
	if err != nil {
		return "", err
	}

	// 生成新的令牌对
	tokenPair, err := s.GenerateTokenPair(user.ID, user.UserAccount)
	if err != nil {
		return "", err
	}

	return tokenPair.AccessToken, nil
}

// StoreRefreshToken 存储刷新令牌到数据库
func (s *TokenService) StoreRefreshToken(userID string, tokenID string, expiresAt time.Time, platform string) error {
	// 计算过期时间
	ttl := expiresAt.Sub(time.Now())

	// RefreshToken存储配置
	ctx := context.Background()
	tokenKey := fmt.Sprintf("refresh_token:%s:%s", userID, tokenID)
	tokenData := map[string]interface{}{
		"userId":     userID,
		"tokenId":    tokenID,
		"expiresAt":  expiresAt.Unix(),
		"created_at": time.Now().Unix(),
		"platform":   platform,
	}
	tokenJson, err := json.Marshal(tokenData)
	if err != nil {
		global.CHAT_LOG.Error("存储RefreshToken序列化失败", "err", err.Error())
		return err
	}
	// RefreshToken存入redis
	if err := global.CHAT_REDIS.Set(ctx, tokenKey, tokenJson, ttl).Err(); err != nil {
		global.CHAT_LOG.Error("存储RefreshToken失败", "err", err.Error())
		return err
	}

	// 存储用户令牌集合
	userTokenKey := fmt.Sprintf("user_tokens:%s", userID)
	if err := global.CHAT_REDIS.SAdd(ctx, userTokenKey, tokenID).Err(); err != nil {
		global.CHAT_LOG.Error("存储RefreshToken失败", "err", err.Error())
		return err
	}
	global.CHAT_REDIS.Expire(ctx, userTokenKey, time.Hour*24*30)

	return nil
}

// 检查令牌是否被撤销
func (s *TokenService) isTokenRevoked(userID string, tokenID string) (bool, error) {
	// 实现检查逻辑
	tokenKey := fmt.Sprintf("refresh_token:%s:%s", userID, tokenID)
	exists, err := global.CHAT_REDIS.Exists(context.Background(), tokenKey).Result()
	if err != nil {
		global.CHAT_LOG.Error("检查RefreshToken是否被撤销，操作失败", "err", err.Error())
		return false, err
	}
	if exists == 0 {
		return true, nil
	}

	return false, nil
}

// 根据ID获取用户
func (s *TokenService) getUserByID(userID string) (*model.User, error) {
	// 实现用户查询逻辑
	queryUser := model.User{}
	err := global.CHAT_MYSQL.Where("id = ?", userID).First(&queryUser).Error
	if err != nil {
		return nil, common.NewServiceError(common.USER_ID_NOT_FOUND)
	}

	return &queryUser, nil
}

// RevokeAllUserTokens 撤销用户的所有令牌（登出所有设备）
func (s *TokenService) RevokeAllUserTokens(userID uint) error {
	// 实现撤销逻辑
	// 例如：将用户所有令牌的isRevoked设为true
	return nil
}

// RevokeToken 撤销特定令牌（单设备登出）
func (s *TokenService) RevokeToken(userID uint, tokenID string) error {
	// 实现撤销逻辑
	return nil
}
