package models

import (
	"time"

	"gorm.io/gorm"
)

// Post 博文模型：对应数据库里的 posts 表
type Post struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Title     string         `gorm:"type:varchar(200);not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	AuthorID  uint           `gorm:"not null;index" json:"author_id"`
	Author    User           `gorm:"foreignKey:AuthorID" json:"author"`
	Category  string         `gorm:"type:varchar(50)" json:"category"`
	Tags      string         `gorm:"type:varchar(200)" json:"tags"`
	Views     uint           `gorm:"default:0" json:"views"`
}
