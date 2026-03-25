package controller

import (
	"blog-system/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreatePostRequest 创建博文请求参数
type CreatePostRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Category string `json:"category"`
	Tags     string `json:"tags"`
}

// UpdatePostRequest 更新博文请求参数
type UpdatePostRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Category string `json:"category"`
	Tags     string `json:"tags"`
}

// CreatePost 创建博文接口
func CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误：标题和内容不能为空")
		return
	}

	// 获取当前登录用户ID
	userID, _ := c.Get("user_id")

	post, err := service.CreatePost(req.Title, req.Content, req.Category, req.Tags, userID.(uint))
	if err != nil {
		Fail(c, "创建博文失败："+err.Error())
		return
	}

	Success(c, post)
}

// GetPostDetail 获取博文详情接口
func GetPostDetail(c *gin.Context) {
	// 从URL路径获取博文ID（如 /posts/1）
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		Fail(c, "博文ID错误")
		return
	}

	post, err := service.GetPostDetail(uint(postID))
	if err != nil {
		Fail(c, err.Error())
		return
	}

	Success(c, post)
}

// GetPostList 获取博文列表接口
func GetPostList(c *gin.Context) {
	// 从URL查询参数获取分页信息（如 /posts?page=1&page_size=10）
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	posts, total, err := service.GetPostList(page, pageSize)
	if err != nil {
		Fail(c, "获取博文列表失败")
		return
	}

	Success(c, gin.H{
		"list":  posts,
		"total": total,
		"page":  page,
	})
}

// UpdatePost 更新博文接口
func UpdatePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		Fail(c, "博文ID错误")
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误：标题和内容不能为空")
		return
	}

	userID, _ := c.Get("user_id")

	post, err := service.UpdatePost(uint(postID), userID.(uint), req.Title, req.Content, req.Category, req.Tags)
	if err != nil {
		Fail(c, err.Error())
		return
	}

	Success(c, post)
}

// DeletePost 删除博文接口
func DeletePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		Fail(c, "博文ID错误")
		return
	}

	userID, _ := c.Get("user_id")

	err = service.DeletePost(uint(postID), userID.(uint))
	if err != nil {
		Fail(c, err.Error())
		return
	}

	Success(c, nil)
}
