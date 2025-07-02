package model

import (
	"time"
)

type Article struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Title        string     `gorm:"type:varchar(255)" json:"title"`
	Content      string     `json:"content"`
	Summary      string     `gorm:"type:varchar(255)" json:"summary"`
	Cover        string     `gorm:"type:varchar(255)" json:"cover"`
	AuthorID     uint       `json:"author_id"`
	Author       User       `gorm:"foreignKey:AuthorID" json:"author"`
	CategoryID   uint       `json:"category_id"`
	Category     Category   `gorm:"foreignKey:CategoryID" json:"category"`
	Tags         []Tag      `gorm:"many2many:article_tags" json:"tags"`
	Status       string     `gorm:"type:varchar(32)" json:"status"` // draft/published/private/pending/rejected
	ViewCount    uint       `json:"view_count"`
	LikeCount    uint       `json:"like_count"`
	CommentCount uint       `json:"comment_count"`
	PublishTime  *time.Time `json:"publish_time"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// 文章-标签关联表
type ArticleTag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ArticleID uint      `json:"article_id"`
	TagID     uint      `json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}
