package service

import (
	"blog-system/dao"
	"blog-system/models"
	"blog-system/utils"
	"errors"
)

// Register 用户注册业务
func Register(username, password, nickname string) error {
	// 1. 检查用户名是否已存在
	_, err := dao.GetUserByUsername(username)
	if err == nil { // 如果没报错，说明用户已存在
		return errors.New("用户名已被注册")
	}

	// 2. 加密密码
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	// 3. 创建用户对象
	user := &models.User{
		Username: username,
		Password: hashedPassword,
		Nickname: nickname,
	}

	// 4. 保存到数据库
	return dao.CreateUser(user)
}

// Login 用户登录业务
func Login(username, password string) (*models.User, string, error) {
	// 1. 根据用户名查询用户
	user, err := dao.GetUserByUsername(username)
	if err != nil {
		return nil, "", errors.New("用户名或密码错误")
	}

	// 2. 验证密码
	if !utils.CheckPassword(password, user.Password) {
		return nil, "", errors.New("用户名或密码错误")
	}

	// 3. 生成JWT Token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GetUserInfo 获取用户信息业务
func GetUserInfo(userID uint) (*models.User, error) {
	return dao.GetUserByID(userID)
}
