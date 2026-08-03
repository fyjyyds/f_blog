package repository

import (
	"f_blog/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// Count 统计评论总数
func (r *CommentRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Comment{}).Count(&count).Error
	return count, err
}

func (r *CommentRepository) Create(tx *gorm.DB, comment *model.Comment) error {
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	return tx.Create(comment).Error
}

func (r *CommentRepository) FindByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.First(&comment, id).Error
	return &comment, err
}

func (r *CommentRepository) FindByArticleID(articleID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Where("article_id = ?", articleID).Order("created_at asc").Find(&comments).Error
	return comments, err
}

func (r *CommentRepository) FindByUserID(userID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&comments).Error
	return comments, err
}

func (r *CommentRepository) Delete(tx *gorm.DB, id uint) error {
	return tx.Delete(&model.Comment{}, id).Error
}

func (r *CommentRepository) AdminList(page, pageSize int, status string) ([]model.Comment, int64, error) {
	db := r.db.Model(&model.Comment{})

	if status != "" {
		db = db.Where("status = ?", status)
	}

	var total int64
	db.Count(&total)

	var comments []model.Comment
	err := db.Preload("User").Preload("Article").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at desc").Find(&comments).Error
	return comments, total, err
}

func (r *CommentRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", id).Update("status", status).Error
}

func (r *CommentRepository) Update(comment *model.Comment, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return r.db.Model(comment).Updates(updates).Error
}
