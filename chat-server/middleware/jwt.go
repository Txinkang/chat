package middleware

import (
	"chat-server/global"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWT 声明结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 不需要验证的路径
var excludePaths = []string{
	"/api/v1/user/register",
	"/api/v1/user/login",
	//"/api/v1/user/test",
	// 可以添加更多不需要验证的路径
}

// JWTAuthWithExclusions 带路径排除的JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前请求路径
		path := c.Request.URL.Path

		// 检查是否在排除路径列表中
		for _, excludePath := range excludePaths {
			if strings.HasPrefix(path, excludePath) {
				// 如果在排除列表中，跳过验证
				c.Next()
				return
			}
		}

		// 不在排除列表中，执行JWT验证
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"data":    nil,
				"message": "未授权",
			})
			c.Abort()
			return
		}

		// 移除Bearer前缀
		if strings.HasPrefix(token, "Bearer ") {
			token = token[7:]
		}

		// 验证token
		claims, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"data":    nil,
				"message": "未授权",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("claims", claims)
		c.Next()
	}
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.CHAT_CONFIG.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username string) (string, error) {
	// 创建声明
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(global.CHAT_CONFIG.JWT.ExpiresTime))), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                                         // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                                                         // 生效时间
			Issuer:    global.CHAT_CONFIG.JWT.Issuer,                                                                          // 签发人
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString([]byte(global.CHAT_CONFIG.JWT.Secret))
}
