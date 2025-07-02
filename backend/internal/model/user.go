package model

import (
	"time"
)

type User struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	Username          string     `gorm:"unique;not null" json:"username"`
	Email             string     `gorm:"unique;not null" json:"email"`
	Password          string     `gorm:"not null" json:"-"`
	Role              string     `gorm:"default:user" json:"role"`
	Status            string     `gorm:"default:active" json:"status"`
	EmailVerified     bool       `gorm:"default:false" json:"email_verified"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	Nickname          string     `json:"nickname"`
	Bio               string     `json:"bio"`
	Gender            string     `json:"gender"` // male/female/other
	Birthday          *time.Time `json:"birthday"`
	Website           string     `json:"website"`
	Github            string     `json:"github"`
	Twitter           string     `json:"twitter"`
	Linkedin          string     `json:"linkedin"`
	Avatar            string     `json:"avatar"`
	EmailVerifyToken  string     `json:"-"` // 邮箱验证token
	EmailVerifyExpire *time.Time `json:"-"` // token过期时间
}
