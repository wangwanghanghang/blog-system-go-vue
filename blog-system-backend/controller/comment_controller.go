package controller

import (
	"blog-system/config"
	"blog-system/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	PostID   uint   `json:"post_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
	ParentID uint   `json:"parent_id"` // 可选，回复评论时使用
}

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Fail(c, "请先登录")
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}

	// 验证文章是否存在且允许评论
	var post models.Post
	if err := config.DB.First(&post, req.PostID).Error; err != nil {
		Fail(c, "文章不存在")
		return
	}

	if !post.AllowComment {
		Fail(c, "该文章不允许评论")
		return
	}

	// 处理父评论ID
	var parentID *uint
	if req.ParentID != 0 {
		parentID = &req.ParentID
	}

	// 创建评论
	comment := models.Comment{
		PostID:    req.PostID,
		UserID:    userID.(uint),
		Content:   req.Content,
		ParentID:  parentID,
		Status:    1,
		IP:        c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		Fail(c, "评论发表失败")
		return
	}

	// 更新文章评论数
	config.DB.Model(&post).UpdateColumn("comment_count", config.DB.Model(&models.Comment{}).Where("post_id = ? AND status = ?", req.PostID, 1).Select("count(*)"))

	// 返回评论信息（包含用户信息）
	config.DB.Preload("User").First(&comment, comment.ID)

	Success(c, comment)
}

// GetPostComments 获取文章评论列表
func GetPostComments(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Fail(c, "文章ID格式错误")
		return
	}

	// 获取分页参数
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if pageInt, err := strconv.Atoi(p); err == nil && pageInt > 0 {
			page = pageInt
		}
	}

	offset := (page - 1) * pageSize

	var comments []models.Comment
	var total int64

	// 只获取顶级评论（parent_id IS NULL）
	config.DB.Model(&models.Comment{}).Where("post_id = ? AND status = ? AND parent_id IS NULL", postID, 1).Count(&total)

	config.DB.Where("post_id = ? AND status = ? AND parent_id IS NULL", postID, 1).
		Preload("User").
		Preload("Replies", "status = ?", 1).
		Preload("Replies.User").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&comments)

	Success(c, gin.H{
		"comments": comments,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Fail(c, "评论ID格式错误")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		Fail(c, "请先登录")
		return
	}

	var comment models.Comment
	if err := config.DB.First(&comment, commentID).Error; err != nil {
		Fail(c, "评论不存在")
		return
	}

	// 检查权限：只有评论作者或管理员可以删除
	var user models.User
	config.DB.First(&user, userID)

	if comment.UserID != userID.(uint) && !user.IsAdmin {
		Fail(c, "权限不足")
		return
	}

	// 软删除评论
	if err := config.DB.Delete(&comment).Error; err != nil {
		Fail(c, "删除评论失败")
		return
	}

	// 更新文章评论数
	config.DB.Model(&models.Post{}).Where("id = ?", comment.PostID).
		UpdateColumn("comment_count", config.DB.Model(&models.Comment{}).Where("post_id = ? AND status = ?", comment.PostID, 1).Select("count(*)"))

	Success(c, nil)
}
