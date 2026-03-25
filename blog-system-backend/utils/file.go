package utils

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"
)

// GenerateFileName 生成随机文件名，防止重名
func GenerateFileName(originalFilename string) string {
	// 获取文件后缀
	ext := filepath.Ext(originalFilename)
	// 生成时间戳 + 随机数的文件名
	timestamp := time.Now().Format("20060102150405")
	randNum := rand.Intn(1000)
	return fmt.Sprintf("%s_%d%s", timestamp, randNum, ext)
}
