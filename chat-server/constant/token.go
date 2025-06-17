package constant

import (
	"chat-server/global"
	"time"
)

var (
	AccessTokenExpireTime  = time.Minute * time.Duration(global.CHAT_CONFIG.JWT.AccessTime)        // 30 minute
	RefreshTokenExpireTime = time.Hour * 24 * time.Duration(global.CHAT_CONFIG.JWT.RefreshTime)    // 30 hour
	UserTokensExpireTime   = time.Hour * 24 * time.Duration(global.CHAT_CONFIG.JWT.UserTokensTime) // 30 hour
	RefreshTokenPrefix     = "refresh_token"
	UserTokensPrefix       = "user_tokens"
)
