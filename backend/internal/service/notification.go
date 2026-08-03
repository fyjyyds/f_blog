package service

import (
	"f_blog/backend/internal/model"
	"f_blog/backend/internal/repository"
)

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) List(userID uint) ([]model.Notification, error) {
	return s.repo.FindByUserID(userID)
}

func (s *NotificationService) Get(id, userID uint) (*model.Notification, error) {
	return s.repo.FindByIDAndUserID(id, userID)
}

func (s *NotificationService) Create(notification *model.Notification) error {
	return s.repo.Create(notification)
}

func (s *NotificationService) MarkRead(id, userID uint) error {
	return s.repo.MarkRead(id, userID)
}

func (s *NotificationService) MarkAllRead(userID uint) error {
	return s.repo.MarkAllRead(userID)
}

func (s *NotificationService) Delete(id, userID uint) error {
	return s.repo.Delete(id, userID)
}
