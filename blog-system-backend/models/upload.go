package models

import (
	"time"

	"gorm.io/gorm"
)

// Upload 文件上传模型
type Upload struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 文件基本信息
	FileName     string `gorm:"type:varchar(255);not null" json:"file_name"`
	OriginalName string `gorm:"type:varchar(255);not null" json:"original_name"`
	FileSize     int64  `gorm:"not null" json:"file_size"`
	MimeType     string `gorm:"type:varchar(100);not null" json:"mime_type"`
	Extension    string `gorm:"type:varchar(20);not null" json:"extension"`

	// 存储信息
	StoragePath string `gorm:"type:varchar(500);not null" json:"storage_path"`
	AccessURL   string `gorm:"type:varchar(500);not null" json:"access_url"`
	FileHash    string `gorm:"type:varchar(64);index" json:"file_hash"` // MD5或SHA1

	// 文件分类
	Category string `gorm:"type:varchar(50);default:'image'" json:"category"` // image, video, document, etc.
	Status   int    `gorm:"default:1;index" json:"status"`                    // 1:正常 0:禁用

	// 关联用户
	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// 使用统计
	UsageCount int `gorm:"default:0" json:"usage_count"` // 被引用次数
}
