package models

import (
	"time"

	"gorm.io/gorm"
)

// Category 分类模型
type Category struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"name"`
	Slug        string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	Color       string         `gorm:"type:varchar(20);default:'primary'" json:"color"`
	Icon        string         `gorm:"type:varchar(50)" json:"icon"`
	Sort        int            `gorm:"default:0;index" json:"sort"`
	Status      int            `gorm:"default:1;index" json:"status"` // 1:启用 0:禁用
	PostCount   int            `gorm:"default:0" json:"post_count"`   // 文章数量统计

	// 关联
	Posts []Post `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
}
