package model

import "time"

type Banner struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"type:varchar(255)" json:"title"`
	Image     string    `gorm:"type:varchar(500)" json:"image"`
	Link      string    `gorm:"type:varchar(500)" json:"link"`
	SortOrder int       `json:"sort_order"`
	Status    string    `gorm:"type:varchar(32)" json:"status"` // active/inactive
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
