package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Status    int            `gorm:"default:1;index" json:"status"` // 1:正常 0:隐藏 -1:删除

	// 关联字段
	PostID   uint  `gorm:"not null;index" json:"post_id"`
	UserID   uint  `gorm:"not null;index" json:"user_id"`
	ParentID *uint `gorm:"index" json:"parent_id"` // 父评论ID，用于回复功能

	// IP和设备信息
	IP        string `gorm:"type:varchar(45)" json:"ip"`
	UserAgent string `gorm:"type:varchar(500)" json:"user_agent"`

	// 关联模型
	Post    Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
	User    User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Parent  *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`

	// 点赞相关
	LikeCount int `gorm:"default:0" json:"like_count"`
}
