package model

import "time"

type Follow struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FollowerID  uint      `json:"follower_id"`  // 关注者ID
	FollowingID uint      `json:"following_id"` // 被关注者ID
	CreatedAt   time.Time `json:"created_at"`
}


