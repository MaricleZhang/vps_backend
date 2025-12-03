package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mariclezhang/vps_backend/internal/util"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.Unauthorized(c, "未授权，请先登录")
			c.Abort()
			return
		}

		// Bearer token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := util.ParseToken(token)
		if err != nil {
			util.Unauthorized(c, "无效的token")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	return userID.(int64), true
}

// GetEmail 从上下文中获取用户邮箱
func GetEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("email")
	if !exists {
		return "", false
	}
	return email.(string), true
}
