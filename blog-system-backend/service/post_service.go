package service

import (
	"blog-system/dao"
	"blog-system/models"
	"errors"
)

// CreatePost 创建博文业务
func CreatePost(title, content, category, tags string, authorID uint) (*models.Post, error) {
	post := &models.Post{
		Title:    title,
		Content:  content,
		AuthorID: authorID,
		Category: category,
		Tags:     tags,
	}
	err := dao.CreatePost(post)
	return post, err
}

// GetPostDetail 获取博文详情业务（同时增加阅读量）
func GetPostDetail(postID uint) (*models.Post, error) {
	// 1. 先查询博文
	post, err := dao.GetPostByID(postID)
	if err != nil {
		return nil, errors.New("博文不存在")
	}

	// 2. 增加阅读量（忽略错误，不影响主流程）
	dao.IncrementPostViews(postID)

	return post, nil
}

// GetPostList 获取博文列表业务
func GetPostList(page, pageSize int, keyword string) ([]models.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // 默认每页10条
	}
	return dao.GetPostList(page, pageSize, keyword)
}

// UpdatePost 更新博文业务（只能更新自己的博文）
func UpdatePost(postID, authorID uint, title, content, category, tags string) (*models.Post, error) {
	// 1. 先查询博文
	post, err := dao.GetPostByID(postID)
	if err != nil {
		return nil, errors.New("博文不存在")
	}

	// 2. 验证是否是作者本人
	if post.AuthorID != authorID {
		return nil, errors.New("无权修改他人的博文")
	}

	// 3. 更新字段
	post.Title = title
	post.Content = content
	post.Category = category
	post.Tags = tags

	// 4. 保存
	err = dao.UpdatePost(post)
	return post, err
}

// DeletePost 删除博文业务（只能删除自己的博文）
func DeletePost(postID, authorID uint) error {
	// 1. 先查询博文
	post, err := dao.GetPostByID(postID)
	if err != nil {
		return errors.New("博文不存在")
	}

	// 2. 验证是否是作者本人
	if post.AuthorID != authorID {
		return errors.New("无权删除他人的博文")
	}

	// 3. 删除
	return dao.DeletePost(postID)
}
