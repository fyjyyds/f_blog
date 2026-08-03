package repository

import (
	"f_blog/backend/internal/model"

	"gorm.io/gorm"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) FirstOrCreate(like *model.Like) error {
	return r.db.Where(like).FirstOrCreate(like).Error
}

func (r *LikeRepository) Delete(userID uint, targetType string, targetID uint) error {
	return r.db.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		Delete(&model.Like{}).Error
}

func (r *LikeRepository) CountByTarget(targetType string, targetID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Like{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Count(&count).Error
	return count, err
}

func (r *LikeRepository) ExistsByUserAndTarget(userID uint, targetType string, targetID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Like{}).
		Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		Count(&count).Error
	return count > 0, err
}

func (r *LikeRepository) UpdateTargetLikeCount(targetType string, targetID uint) error {
	count, err := r.CountByTarget(targetType, targetID)
	if err != nil {
		return err
	}
	switch targetType {
	case "article":
		return r.db.Model(&model.Article{}).Where("id = ?", targetID).Update("like_count", count).Error
	case "comment":
		return r.db.Model(&model.Comment{}).Where("id = ?", targetID).Update("like_count", count).Error
	}
	return nil
}

func (r *LikeRepository) FindByUserID(userID uint, targetType string) ([]model.Like, error) {
	var likes []model.Like
	err := r.db.Where("user_id = ? AND target_type = ?", userID, targetType).
		Order("created_at desc").Find(&likes).Error
	return likes, err
}
