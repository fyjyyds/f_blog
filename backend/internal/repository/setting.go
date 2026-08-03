package repository

import (
	"f_blog/backend/internal/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SettingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

func (r *SettingRepository) FindAll() ([]model.Setting, error) {
	var settings []model.Setting
	err := r.db.Find(&settings).Error
	return settings, err
}

func (r *SettingRepository) Upsert(key string, value interface{}, valueType string) error {
	var setting model.Setting
	result := r.db.Where("`key` = ?", key).First(&setting)
	if result.Error != nil {
		setting = model.Setting{
			Key:   key,
			Type:  valueType,
			Value: fmt.Sprintf("%v", value),
		}
		return r.db.Create(&setting).Error
	}
	setting.Value = fmt.Sprintf("%v", value)
	setting.UpdatedAt = time.Now()
	return r.db.Save(&setting).Error
}

func (r *SettingRepository) ResetAll() error {
	return r.db.Where("1 = 1").Delete(&model.Setting{}).Error
}
