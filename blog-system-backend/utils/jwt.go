package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 定义JWT的密钥（自己随便改一个复杂点的字符串，不要泄露）
var jwtSecret = []byte("htyWCIHh3hIipPtlQk/ztxuQ+vU7jLysJr2Rh/T/J1M=")

// Claims 自定义JWT声明结构：存用户ID和用户名，用于识别用户身份
type Claims struct {
	UserID               uint   `json:"user_id"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        // 嵌入JWT标准声明（过期时间等）
}

// GenerateToken 生成JWT Token：用户登录成功后返回给前端，前端后续请求带这个Token证明身份
func GenerateToken(userID uint, username string) (string, error) {
	// 设置Token过期时间：7天
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	// 创建声明
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // 签发时间
			Issuer:    "blog-system",                      // 签发者
		},
	}

	// 用HS256算法签名生成Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析JWT Token：验证前端传来的Token是否有效，并提取用户信息
func ParseToken(tokenString string) (*Claims, error) {
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证Token是否有效
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}