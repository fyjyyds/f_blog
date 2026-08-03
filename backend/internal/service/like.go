package service

import (
	"f_blog/backend/internal/model"
	"f_blog/backend/internal/repository"
	"fmt"
	"log"
)

type LikeService struct {
	likeRepo  *repository.LikeRepository
	notifRepo *repository.NotificationRepository
	articleRepo *repository.ArticleRepository
	commentRepo *repository.CommentRepository
}

func NewLikeService(
	likeRepo *repository.LikeRepository,
	notifRepo *repository.NotificationRepository,
	articleRepo *repository.ArticleRepository,
	commentRepo *repository.CommentRepository,
) *LikeService {
	return &LikeService{
		likeRepo:    likeRepo,
		notifRepo:   notifRepo,
		articleRepo: articleRepo,
		commentRepo: commentRepo,
	}
}

type LikeResult struct {
	Liked bool  `json:"liked"`
	Count int64 `json:"count"`
}

func (s *LikeService) Like(userID uint, targetType string, targetID uint) (*LikeResult, error) {
	like := &model.Like{
		UserID:     userID,
		TargetType: targetType,
		TargetID:   targetID,
	}
	if err := s.likeRepo.FirstOrCreate(like); err != nil {
		return nil, err
	}

	count, _ := s.likeRepo.CountByTarget(targetType, targetID)
	_ = s.likeRepo.UpdateTargetLikeCount(targetType, targetID)

	// 发送通知
	if targetType == "article" {
		article, err := s.articleRepo.FindByID(targetID)
		if err == nil && article.AuthorID != userID {
			if err := s.notifRepo.Create(&model.Notification{
				UserID:  article.AuthorID,
				Type:    "like",
				Title:   "你的文章被点赞",
				Content: "有人点赞了你的文章：" + article.Title,
				Data:    fmt.Sprintf(`{"article_id":%d}`, article.ID),
			}); err != nil {
				log.Printf("发送点赞通知失败: %v", err)
			}
		}
	} else if targetType == "comment" {
		comment, err := s.commentRepo.FindByID(targetID)
		if err == nil && comment.UserID != userID {
			if err := s.notifRepo.Create(&model.Notification{
				UserID:  comment.UserID,
				Type:    "like",
				Title:   "你的评论被点赞",
				Content: "有人点赞了你的评论",
				Data:    fmt.Sprintf(`{"comment_id":%d}`, comment.ID),
			}); err != nil {
				log.Printf("发送点赞通知失败: %v", err)
			}
		}
	}

	return &LikeResult{Liked: true, Count: count}, nil
}

func (s *LikeService) Unlike(userID uint, targetType string, targetID uint) (*LikeResult, error) {
	if err := s.likeRepo.Delete(userID, targetType, targetID); err != nil {
		return nil, err
	}

	count, _ := s.likeRepo.CountByTarget(targetType, targetID)
	_ = s.likeRepo.UpdateTargetLikeCount(targetType, targetID)

	return &LikeResult{Liked: false, Count: count}, nil
}

func (s *LikeService) Status(userID uint, targetType string, targetID uint) (*LikeResult, error) {
	count, err := s.likeRepo.CountByTarget(targetType, targetID)
	if err != nil {
		return nil, err
	}

	liked := false
	if userID != 0 {
		liked, _ = s.likeRepo.ExistsByUserAndTarget(userID, targetType, targetID)
	}

	return &LikeResult{Liked: liked, Count: count}, nil
}

func (s *LikeService) ListByUserID(userID uint, targetType string) ([]model.Like, error) {
	return s.likeRepo.FindByUserID(userID, targetType)
}
