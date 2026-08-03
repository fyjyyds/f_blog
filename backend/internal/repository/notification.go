package repository

import (
	"f_blog/backend/internal/model"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) FindByUserID(userID uint) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) FindByIDAndUserID(id, userID uint) (*model.Notification, error) {
	var notification model.Notification
	err := r.db.Where("user_id = ?", userID).First(&notification, id).Error
	return &notification, err
}

func (r *NotificationRepository) Create(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepository) MarkRead(id, userID uint) error {
	return r.db.Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

func (r *NotificationRepository) MarkAllRead(userID uint) error {
	return r.db.Model(&model.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error
}

func (r *NotificationRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Notification{}).Error
}
