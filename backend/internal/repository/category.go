package repository

import (
	"f_blog/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindAll() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Order("sort_order asc, created_at desc").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) FindByID(id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, id).Error
	return &category, err
}

func (r *CategoryRepository) Create(category *model.Category) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	return r.db.Create(category).Error
}

func (r *CategoryRepository) Update(category *model.Category, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return r.db.Model(category).Updates(updates).Error
}

func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}
