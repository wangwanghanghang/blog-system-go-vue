package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型：对应数据库里的 users 表
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"` // 主键ID，前端用id接收
	CreatedAt time.Time      `json:"created_at"`           // 创建时间，前端用created_at接收
	UpdatedAt time.Time      `json:"updated_at"`           // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`       // 软删除字段，不返回给前端
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Nickname  string         `gorm:"type:varchar(50)" json:"nickname"`
}
