package controller

import (
	"blog-system/config"
	"blog-system/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCategories 获取分类列表
func GetCategories(c *gin.Context) {
	var categories []models.Category

	result := config.DB.Where("status = ?", 1).Order("sort ASC, id ASC").Find(&categories)
	if result.Error != nil {
		Fail(c, "获取分类列表失败")
		return
	}

	Success(c, categories)
}

// GetCategoryDetail 获取分类详情（包含文章列表）
func GetCategoryDetail(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Fail(c, "分类ID格式错误")
		return
	}

	var category models.Category
	result := config.DB.Where("id = ? AND status = ?", categoryID, 1).First(&category)
	if result.Error != nil {
		Fail(c, "分类不存在")
		return
	}

	// 获取分页参数
	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		if pageInt, err := strconv.Atoi(p); err == nil && pageInt > 0 {
			page = pageInt
		}
	}

	// 获取该分类下的文章
	var posts []models.Post
	var total int64

	offset := (page - 1) * pageSize

	config.DB.Model(&models.Post{}).Where("category_id = ? AND status = ?", categoryID, 1).Count(&total)
	config.DB.Where("category_id = ? AND status = ?", categoryID, 1).
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&posts)

	Success(c, gin.H{
		"category": category,
		"posts":    posts,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// CreateCategory 创建分类（管理员功能）
func CreateCategory(c *gin.Context) {
	// 检查管理员权限
	userID, exists := c.Get("user_id")
	if !exists {
		Fail(c, "未登录")
		return
	}

	var user models.User
	config.DB.First(&user, userID)
	if !user.IsAdmin {
		Fail(c, "权限不足")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Slug        string `json:"slug" binding:"required"`
		Description string `json:"description"`
		Color       string `json:"color"`
		Icon        string `json:"icon"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}

	category := models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		Status:      1,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		Fail(c, "创建分类失败")
		return
	}

	Success(c, category)
}
