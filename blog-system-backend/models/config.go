package models

import (
	"time"

	"gorm.io/gorm"
)

// Config 系统配置模型
type Config struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Key         string `gorm:"type:varchar(100);not null;uniqueIndex" json:"key"`
	Value       string `gorm:"type:text" json:"value"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	Type        string `gorm:"type:varchar(20);default:'string'" json:"type"` // string, int, bool, json
	Group       string `gorm:"type:varchar(50);default:'system'" json:"group"`
	Sort        int    `gorm:"default:0" json:"sort"`
	IsPublic    bool   `gorm:"default:false" json:"is_public"` // 是否对外公开
}
