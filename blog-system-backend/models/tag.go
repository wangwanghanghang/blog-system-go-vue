package models

import (
	"time"

	"gorm.io/gorm"
)

// Tag 标签模型
type Tag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"name"`
	Slug      string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"slug"`
	Color     string         `gorm:"type:varchar(20);default:'default'" json:"color"`
	Sort      int            `gorm:"default:0;index" json:"sort"`
	Status    int            `gorm:"default:1;index" json:"status"` // 1:启用 0:禁用
	PostCount int            `gorm:"default:0" json:"post_count"`   // 文章数量统计

	// 多对多关联
	Posts []Post `gorm:"many2many:post_tags" json:"posts,omitempty"`
}

// PostTag 文章标签关联表
type PostTag struct {
	PostID uint `gorm:"primaryKey" json:"post_id"`
	TagID  uint `gorm:"primaryKey" json:"tag_id"`
	Post   Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Tag    Tag  `gorm:"foreignKey:TagID" json:"tag,omitempty"`
}
