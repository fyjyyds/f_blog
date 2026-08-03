package service

import (
	"errors"
	"f_blog/backend/internal/model"
	"f_blog/backend/internal/repository"
	"fmt"
	"log"
)

type FollowService struct {
	followRepo *repository.FollowRepository
	userRepo   *repository.UserRepository
	notifRepo  *repository.NotificationRepository
}

func NewFollowService(
	followRepo *repository.FollowRepository,
	userRepo *repository.UserRepository,
	notifRepo *repository.NotificationRepository,
) *FollowService {
	return &FollowService{
		followRepo: followRepo,
		userRepo:   userRepo,
		notifRepo:  notifRepo,
	}
}

type FollowWithUser struct {
	model.Follow
	Username string `json:"following_username,omitempty"`
	Nickname string `json:"following_nickname,omitempty"`
}

type FollowerWithUser struct {
	model.Follow
	Username string `json:"follower_username,omitempty"`
	Nickname string `json:"follower_nickname,omitempty"`
}

func (s *FollowService) Follow(userID, followingID uint) error {
	if userID == followingID {
		return errors.New("不能关注自己")
	}
	follow := &model.Follow{
		FollowerID:  userID,
		FollowingID: followingID,
	}
	if err := s.followRepo.FirstOrCreate(follow); err != nil {
		return err
	}

	// 发送通知
	if err := s.notifRepo.Create(&model.Notification{
		UserID:  followingID,
		Type:    "follow",
		Title:   "你有新粉丝",
		Content: "有人关注了你",
		Data:    fmt.Sprintf(`{"follower_id":%d}`, userID),
	}); err != nil {
		log.Printf("发送关注通知失败: %v", err)
	}
	return nil
}

func (s *FollowService) Unfollow(userID, followingID uint) error {
	if userID == followingID {
		return errors.New("不能取关自己")
	}
	return s.followRepo.Delete(userID, followingID)
}

func (s *FollowService) ListFollowings(userID uint) ([]FollowWithUser, error) {
	follows, err := s.followRepo.FindFollowings(userID)
	if err != nil {
		return nil, err
	}

	ids := make([]uint, 0, len(follows))
	for _, f := range follows {
		ids = append(ids, f.FollowingID)
	}
	users, _ := s.userRepo.FindByIDs(ids)
	userMap := make(map[uint]model.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	result := make([]FollowWithUser, 0, len(follows))
	for _, f := range follows {
		u := userMap[f.FollowingID]
		result = append(result, FollowWithUser{
			Follow:   f,
			Username: u.Username,
			Nickname: u.Nickname,
		})
	}
	return result, nil
}

func (s *FollowService) ListFollowers(userID uint) ([]FollowerWithUser, error) {
	follows, err := s.followRepo.FindFollowers(userID)
	if err != nil {
		return nil, err
	}

	ids := make([]uint, 0, len(follows))
	for _, f := range follows {
		ids = append(ids, f.FollowerID)
	}
	users, _ := s.userRepo.FindByIDs(ids)
	userMap := make(map[uint]model.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	result := make([]FollowerWithUser, 0, len(follows))
	for _, f := range follows {
		u := userMap[f.FollowerID]
		result = append(result, FollowerWithUser{
			Follow:   f,
			Username: u.Username,
			Nickname: u.Nickname,
		})
	}
	return result, nil
}
