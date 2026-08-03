package service

import (
	"errors"
	"f_blog/backend/internal/model"
	"f_blog/backend/internal/repository"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type CommentService struct {
	db          *gorm.DB
	commentRepo *repository.CommentRepository
	articleRepo *repository.ArticleRepository
	notifRepo   *repository.NotificationRepository
	userRepo    *repository.UserRepository
}

func NewCommentService(
	db *gorm.DB,
	commentRepo *repository.CommentRepository,
	articleRepo *repository.ArticleRepository,
	notifRepo *repository.NotificationRepository,
	userRepo *repository.UserRepository,
) *CommentService {
	return &CommentService{
		db:          db,
		commentRepo: commentRepo,
		articleRepo: articleRepo,
		notifRepo:   notifRepo,
		userRepo:    userRepo,
	}
}

func (s *CommentService) Count() (int64, error) {
	return s.commentRepo.Count()
}

type CommentWithDisplayName struct {
	model.Comment
	DisplayName string `json:"display_name"`
}

func (s *CommentService) Create(articleID uint, userID uint, content string, parentID uint, replyToUser string) (*model.Comment, error) {
	comment := &model.Comment{
		ArticleID:   articleID,
		UserID:      userID,
		ParentID:    parentID,
		Content:     content,
		Status:      "approved",
		ReplyToUser: replyToUser,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.commentRepo.Create(tx, comment); err != nil {
			return err
		}
		if err := tx.Model(&model.Article{}).Where("id = ?", articleID).
			Update("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 发送通知给文章作者
	article, err := s.articleRepo.FindByID(articleID)
	if err == nil && article.AuthorID != userID {
		if err := s.notifRepo.Create(&model.Notification{
			UserID:  article.AuthorID,
			Type:    "comment",
			Title:   "你的文章有新评论",
			Content: "有人评论了你的文章：" + article.Title,
			Data:    fmt.Sprintf(`{"article_id":%d}`, article.ID),
		}); err != nil {
			log.Printf("发送评论通知失败: %v", err)
		}
	}

	return comment, nil
}

func (s *CommentService) ListByArticleID(articleID uint) ([]CommentWithDisplayName, error) {
	comments, err := s.commentRepo.FindByArticleID(articleID)
	if err != nil {
		return nil, err
	}

	userIDs := make([]uint, 0, len(comments))
	userIDSet := make(map[uint]struct{})
	for _, c := range comments {
		if _, ok := userIDSet[c.UserID]; !ok {
			userIDs = append(userIDs, c.UserID)
			userIDSet[c.UserID] = struct{}{}
		}
	}
	users, _ := s.userRepo.FindByIDs(userIDs)
	userMap := make(map[uint]model.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	result := make([]CommentWithDisplayName, 0, len(comments))
	for _, c := range comments {
		user := userMap[c.UserID]
		displayName := user.Nickname
		if displayName == "" {
			displayName = user.Username
		}
		result = append(result, CommentWithDisplayName{
			Comment:     c,
			DisplayName: displayName,
		})
	}
	return result, nil
}

func (s *CommentService) ListByUserID(userID uint) ([]model.Comment, error) {
	return s.commentRepo.FindByUserID(userID)
}

func (s *CommentService) Delete(id, userID uint, role string) error {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		return errors.New("评论不存在")
	}
	if comment.UserID != userID && role != "admin" {
		return errors.New("无权限操作")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.commentRepo.Delete(tx, id); err != nil {
			return err
		}
		return tx.Model(&model.Article{}).Where("id = ?", comment.ArticleID).
			Update("comment_count", gorm.Expr("GREATEST(CAST(comment_count AS SIGNED) - 1, 0)")).Error
	})
}

func (s *CommentService) AdminList(page, pageSize int, status string) ([]model.Comment, int64, error) {
	return s.commentRepo.AdminList(page, pageSize, status)
}

func (s *CommentService) AdminApprove(id uint) error {
	return s.commentRepo.UpdateStatus(id, "approved")
}

func (s *CommentService) AdminReject(id uint) error {
	return s.commentRepo.UpdateStatus(id, "rejected")
}

func (s *CommentService) AdminUpdate(id uint, content, status string) (*model.Comment, error) {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("评论不存在")
	}
	updates := map[string]interface{}{
		"content": content,
		"status":  status,
	}
	err = s.commentRepo.Update(comment, updates)
	return comment, err
}

func (s *CommentService) AdminDelete(id uint) error {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		return errors.New("评论不存在")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.commentRepo.Delete(tx, id); err != nil {
			return err
		}
		return tx.Model(&model.Article{}).Where("id = ?", comment.ArticleID).
			Update("comment_count", gorm.Expr("GREATEST(CAST(comment_count AS SIGNED) - 1, 0)")).Error
	})
}
