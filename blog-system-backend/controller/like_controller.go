package controller

import (
	"blog-system/config"
	"blog-system/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ToggleLike 切换点赞状态
func ToggleLike(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Fail(c, "请先登录")
		return
	}

	postID, err := strconv.ParseUint(c.Param("postId"), 10, 32)
	if err != nil {
		Fail(c, "文章ID格式错误")
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		Fail(c, "文章不存在")
		return
	}

	// 检查是否已经点赞
	var existingLike models.Like
	result := config.DB.Where("user_id = ? AND post_id = ? AND type = ?", userID, postID, "post").First(&existingLike)

	if result.Error != nil {
		// 还未点赞，创建点赞记录
		like := models.Like{
			UserID: userID.(uint),
			PostID: uint(postID),
			Type:   "post",
		}

		if err := config.DB.Create(&like).Error; err != nil {
			Fail(c, "点赞失败")
			return
		}

		// 更新文章点赞数
		config.DB.Model(&post).UpdateColumn("like_count", config.DB.Model(&models.Like{}).Where("post_id = ? AND type = ?", postID, "post").Select("count(*)"))

		Success(c, gin.H{
			"action":  "liked",
			"message": "点赞成功",
		})
	} else {
		// 已经点赞，取消点赞
		if err := config.DB.Delete(&existingLike).Error; err != nil {
			Fail(c, "取消点赞失败")
			return
		}

		// 更新文章点赞数
		config.DB.Model(&post).UpdateColumn("like_count", config.DB.Model(&models.Like{}).Where("post_id = ? AND type = ?", postID, "post").Select("count(*)"))

		Success(c, gin.H{
			"action":  "unliked",
			"message": "取消点赞成功",
		})
	}
}

// GetPostLikes 获取文章点赞状态
func GetPostLikes(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 32)
	if err != nil {
		Fail(c, "文章ID格式错误")
		return
	}

	// 获取点赞总数
	var likeCount int64
	config.DB.Model(&models.Like{}).Where("post_id = ? AND type = ?", postID, "post").Count(&likeCount)

	result := gin.H{
		"post_id":    postID,
		"like_count": likeCount,
		"is_liked":   false,
	}

	// 如果用户已登录，检查是否已点赞
	if userID, exists := c.Get("user_id"); exists {
		var existingLike models.Like
		if err := config.DB.Where("user_id = ? AND post_id = ? AND type = ?", userID, postID, "post").First(&existingLike).Error; err == nil {
			result["is_liked"] = true
		}
	}

	Success(c, result)
}
