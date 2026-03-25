package dao

import (
	"blog-system/config"
	"blog-system/models"

	"gorm.io/gorm"
)

// CreatePost 创建博文
func CreatePost(post *models.Post) error {
	return config.DB.Create(post).Error
}

// GetPostByID 根据ID查询博文（带作者信息）
func GetPostByID(postID uint) (*models.Post, error) {
	var post models.Post
	// Preload("Author") 会自动关联查询作者信息
	err := config.DB.Preload("Author").First(&post, postID).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostList 分页查询博文列表（带作者信息，按创建时间倒序）
func GetPostList(page, pageSize int, keyword string) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// 创建查询会话
	query := config.DB.Model(&models.Post{})

	// 如果有搜索关键词
	if keyword != "" {
		likeArgs := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR content LIKE ?", likeArgs, likeArgs)
	}

	// 先查询总条数
	query.Count(&total)

	// 再分页查询数据
	offset := (page - 1) * pageSize
	err := query.Preload("Author").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&posts).Error

	return posts, total, err
}

// UpdatePost 更新博文
func UpdatePost(post *models.Post) error {
	return config.DB.Save(post).Error
}

// DeletePost 删除博文（软删除，因为GORM默认软删除）
func DeletePost(postID uint) error {
	return config.DB.Delete(&models.Post{}, postID).Error
}

// IncrementPostViews 增加博文阅读量
func IncrementPostViews(postID uint) error {
	return config.DB.Model(&models.Post{}).Where("id = ?", postID).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}
