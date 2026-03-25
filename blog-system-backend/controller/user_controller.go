package controller

import (
	"blog-system/service"

	"github.com/gin-gonic/gin"
)

// RegisterRequest 注册请求参数结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required"` // 用户名：必填
	Password string `json:"password" binding:"required"` // 密码：必填
	Nickname string `json:"nickname"`                    // 昵称：选填
}

// LoginRequest 登录请求参数结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册接口
func Register(c *gin.Context) {
	var req RegisterRequest
	// 1. 绑定前端传来的JSON参数
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误：请检查用户名和密码是否填写")
		return
	}

	// 2. 调用Service注册
	err := service.Register(req.Username, req.Password, req.Nickname)
	if err != nil {
		Fail(c, err.Error())
		return
	}

	Success(c, nil)
}

// Login 用户登录接口
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误：请检查用户名和密码是否填写")
		return
	}

	user, token, err := service.Login(req.Username, req.Password)
	if err != nil {
		Fail(c, err.Error())
		return
	}

	// 返回用户信息和Token
	Success(c, gin.H{
		"token":    token,
		"user_id":  user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
	})
}

// GetUserInfo 获取当前登录用户信息接口
func GetUserInfo(c *gin.Context) {
	// 从中间件存入的上下文获取用户ID
	userID, _ := c.Get("user_id")

	user, err := service.GetUserInfo(userID.(uint))
	if err != nil {
		Fail(c, "获取用户信息失败")
		return
	}

	Success(c, user)
}
