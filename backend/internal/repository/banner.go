package repository

import (
	"f_blog/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type BannerRepository struct {
	db *gorm.DB
}

func NewBannerRepository(db *gorm.DB) *BannerRepository {
	return &BannerRepository{db: db}
}

func (r *BannerRepository) FindActive() ([]model.Banner, error) {
	var banners []model.Banner
	err := r.db.Where("status = ?", "active").Order("sort_order asc, created_at desc").Find(&banners).Error
	return banners, err
}

func (r *BannerRepository) FindAll() ([]model.Banner, error) {
	var banners []model.Banner
	err := r.db.Order("sort_order asc, created_at desc").Find(&banners).Error
	return banners, err
}

func (r *BannerRepository) FindByID(id uint) (*model.Banner, error) {
	var banner model.Banner
	err := r.db.First(&banner, id).Error
	return &banner, err
}

func (r *BannerRepository) Create(banner *model.Banner) error {
	banner.CreatedAt = time.Now()
	banner.UpdatedAt = time.Now()
	return r.db.Create(banner).Error
}

func (r *BannerRepository) Update(banner *model.Banner, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return r.db.Model(banner).Updates(updates).Error
}

func (r *BannerRepository) Delete(id uint) error {
	return r.db.Delete(&model.Banner{}, id).Error
}
