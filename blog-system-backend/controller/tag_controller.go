package controller

import (
	"blog-system/config"
	"blog-system/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetTags 获取标签列表
func GetTags(c *gin.Context) {
	var tags []models.Tag

	result := config.DB.Order("post_count DESC, id ASC").Find(&tags)
	if result.Error != nil {
		Fail(c, "获取标签列表失败")
		return
	}

	Success(c, tags)
}

// GetTagDetail 获取标签详情（包含文章列表）
func GetTagDetail(c *gin.Context) {
	tagID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Fail(c, "标签ID格式错误")
		return
	}

	var tag models.Tag
	result := config.DB.First(&tag, tagID)
	if result.Error != nil {
		Fail(c, "标签不存在")
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

	// 通过多对多关联获取该标签下的文章
	var posts []models.Post
	var total int64

	offset := (page - 1) * pageSize

	// 计算总数
	total = config.DB.Model(&tag).Association("Posts").Count()

	// 获取文章列表
	config.DB.Model(&tag).
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Where("status = ?", 1).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Association("Posts").
		Find(&posts)

	Success(c, gin.H{
		"tag":      tag,
		"posts":    posts,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Slug  string `json:"slug" binding:"required"`
		Color string `json:"color"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}

	tag := models.Tag{
		Name:  req.Name,
		Slug:  req.Slug,
		Color: req.Color,
	}

	if err := config.DB.Create(&tag).Error; err != nil {
		Fail(c, "创建标签失败")
		return
	}

	Success(c, tag)
}
