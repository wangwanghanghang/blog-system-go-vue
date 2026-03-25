package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 加密密码：把明文密码转成密文存储，防止数据库泄露后密码被盗
func HashPassword(password string) (string, error) {
	// bcrypt.GenerateFromPassword会自动生成盐值，加密强度设为14（平衡安全和性能）
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword 验证密码：对比明文密码和数据库里的密文是否一致
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // 如果err为nil，说明密码正确
}
