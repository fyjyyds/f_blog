package model

import "time"

type Comment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ArticleID   uint      `json:"article_id"`
	UserID      uint      `json:"user_id"`
	ParentID    uint      `json:"parent_id"`
	Content     string    `gorm:"type:text" json:"content"`
	Status      string    `gorm:"type:varchar(32)" json:"status"`
	LikeCount   uint      `json:"like_count"`
	ReplyCount  uint      `json:"reply_count"`
	ReplyToUser string    `json:"reply_to_user"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	Article     Article   `gorm:"foreignKey:ArticleID" json:"article"`
}
