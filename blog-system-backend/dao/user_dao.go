package dao

import (
	"blog-system/config"
	"blog-system/models"
)

// CreateUser 创建用户：注册时用
func CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}

// GetUserByUsername 根据用户名查询用户：登录时验证用户名是否存在
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据用户ID查询用户：获取用户信息时用
func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
