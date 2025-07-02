package model

import "time"

type Like struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `json:"user_id"`
	TargetType string    `gorm:"type:varchar(32)" json:"target_type"` // "article" or "comment"
	TargetID   uint      `json:"target_id"`
	CreatedAt  time.Time `json:"created_at"`
}
