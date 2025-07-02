package model

import (
	"time"
)

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Type      string    `gorm:"type:enum('comment','like','follow','system','review')" json:"type"`
	Title     string    `gorm:"type:varchar(100)" json:"title"`
	Content   string    `json:"content"`
	Data      string    `json:"data"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
