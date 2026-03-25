package controller

import (
	"blog-system/utils"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// UploadImage 上传图片接口
func UploadImage(c *gin.Context) {
	// 1. 获取前端上传的文件
	file, err := c.FormFile("image")
	if err != nil {
		Fail(c, "请选择要上传的图片")
		return
	}

	// 2. 限制文件大小（最大 5MB）
	if file.Size > 5*1024*1024 {
		Fail(c, "图片大小不能超过5MB")
		return
	}

	// 3. 限制文件类型（只允许图片，处理大小写问题）
	ext := strings.ToLower(filepath.Ext(file.Filename)) // ⭐ 新增：把后缀名转成小写
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		Fail(c, "只支持上传 jpg、jpeg、png、gif、webp 格式的图片")
		return
	}

	// 4. 生成新的文件名
	newFileName := utils.GenerateFileName(file.Filename)
	// 定义保存路径
	savePath := "./uploads/" + newFileName

	// 5. 保存文件到本地
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		Fail(c, "图片保存失败："+err.Error())
		return
	}

	// 6. 返回图片的访问URL
	imageUrl := fmt.Sprintf("http://127.0.0.1:8080/uploads/%s", newFileName)
	Success(c, gin.H{
		"url": imageUrl,
	})
}
