package models

import (
	"time"

	"gorm.io/gorm"
)

// Like 点赞模型
type Like struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint   `gorm:"not null;index" json:"user_id"`
	PostID uint   `gorm:"not null;index" json:"post_id"`               // 统一使用 PostID 存储目标ID，通过 Type 区分
	Type   string `gorm:"type:varchar(20);default:'post'" json:"type"` // post, comment

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	// 注意：这里 Post 关联仅在该 Like真的是针对 Post 时有效。
	// 如果是 Comment，则 PostID 存的是 CommentID，此时关联 Post 会失效或错误。
	// 更好的做法是使用多态关联，或者明确字段 CommentID。
	// 这里为了兼容性(并未看到具体业务代码)，暂保持原样，但在 models 中建议改进。
}
