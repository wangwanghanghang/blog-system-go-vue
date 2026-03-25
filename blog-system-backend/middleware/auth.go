package middleware

import (
	"blog-system/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired JWT鉴权中间件：保护需要登录才能访问的接口
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未登录，请先登录"})
			c.Abort() // 终止请求，不继续执行后面的接口
			return
		}

		// 2. 检查Token格式（格式应该是 "Bearer xxx"）
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token格式错误"})
			c.Abort()
			return
		}

		// 3. 解析Token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token无效或已过期"})
			c.Abort()
			return
		}

		// 4. 把用户信息存入上下文，后面的接口可以通过 c.Get("user_id") 获取
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		// 5. 继续执行后面的接口
		c.Next()
	}
}
