package repository

import (
	"f_blog/backend/internal/model"

	"gorm.io/gorm"
)

type FollowRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

func (r *FollowRepository) FirstOrCreate(follow *model.Follow) error {
	return r.db.Where(follow).FirstOrCreate(follow).Error
}

func (r *FollowRepository) Delete(followerID, followingID uint) error {
	return r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&model.Follow{}).Error
}

func (r *FollowRepository) FindFollowings(userID uint) ([]model.Follow, error) {
	var follows []model.Follow
	err := r.db.Where("follower_id = ?", userID).Find(&follows).Error
	return follows, err
}

func (r *FollowRepository) FindFollowers(userID uint) ([]model.Follow, error) {
	var follows []model.Follow
	err := r.db.Where("following_id = ?", userID).Find(&follows).Error
	return follows, err
}
