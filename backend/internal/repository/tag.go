package repository

import (
	"f_blog/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) FindAll() ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Order("created_at desc").Find(&tags).Error
	return tags, err
}

func (r *TagRepository) FindByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.First(&tag, id).Error
	return &tag, err
}

func (r *TagRepository) Create(tag *model.Tag) error {
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()
	return r.db.Create(tag).Error
}

func (r *TagRepository) Update(tag *model.Tag, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return r.db.Model(tag).Updates(updates).Error
}

func (r *TagRepository) Delete(id uint) error {
	return r.db.Delete(&model.Tag{}, id).Error
}
