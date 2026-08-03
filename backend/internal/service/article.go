package service

import (
	"encoding/json"
	"errors"
	"f_blog/backend/internal/cache"
	"f_blog/backend/internal/model"
	"f_blog/backend/internal/repository"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type ArticleService struct {
	db          *gorm.DB
	articleRepo *repository.ArticleRepository
	notifRepo   *repository.NotificationRepository
}

func NewArticleService(
	db *gorm.DB,
	articleRepo *repository.ArticleRepository,
	notifRepo *repository.NotificationRepository,
) *ArticleService {
	return &ArticleService{
		db:          db,
		articleRepo: articleRepo,
		notifRepo:   notifRepo,
	}
}

type CreateArticleInput struct {
	Title      string
	Content    string
	Summary    string
	Cover      string
	CategoryID uint
	Status     string
	TagIDs     []uint
	AuthorID   uint
	Role       string
}

func (s *ArticleService) Count() (int64, error) {
	return s.articleRepo.Count()
}

func (s *ArticleService) CountByStatus(status string) (int64, error) {
	return s.articleRepo.CountByStatus(status)
}

func (s *ArticleService) Create(input CreateArticleInput) (*model.Article, error) {
	status := "pending"
	if input.Role == "admin" {
		if input.Status != "" {
			status = input.Status
		} else {
			status = "published"
		}
	}

	article := &model.Article{
		Title:      input.Title,
		Content:    input.Content,
		Summary:    input.Summary,
		Cover:      input.Cover,
		AuthorID:   input.AuthorID,
		CategoryID: input.CategoryID,
		Status:     status,
	}
	if status == "published" {
		now := time.Now()
		article.PublishTime = &now
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.articleRepo.Create(tx, article); err != nil {
			return err
		}
		if len(input.TagIDs) > 0 {
			if err := s.articleRepo.ReplaceTags(tx, article.ID, input.TagIDs); err != nil {
				return err
			}
		}
		return nil
	})

	if err == nil {
		s.invalidateListCaches()
	}

	return article, err
}

type UpdateArticleInput struct {
	Title      string
	Content    string
	Summary    string
	Cover      string
	CategoryID uint
	Status     string
	TagIDs     []uint
}

func (s *ArticleService) Update(id uint, input UpdateArticleInput, userID uint, role string) (*model.Article, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("文章不存在")
	}
	if article.AuthorID != userID && role != "admin" {
		return nil, errors.New("无权限操作")
	}

	updates := map[string]interface{}{
		"title":   input.Title,
		"content": input.Content,
		"summary": input.Summary,
		"cover":   input.Cover,
	}
	if input.CategoryID != 0 {
		updates["category_id"] = input.CategoryID
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	if input.Status == "published" && article.PublishTime == nil {
		now := time.Now()
		updates["publish_time"] = &now
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.articleRepo.UpdateWithTx(tx, article, updates); err != nil {
			return err
		}
		if err := s.articleRepo.ReplaceTags(tx, article.ID, input.TagIDs); err != nil {
			return err
		}
		return nil
	})

	if err == nil {
		s.invalidateArticleCaches(id)
	}

	return article, err
}

func (s *ArticleService) Delete(id, userID uint, role string) error {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return errors.New("文章不存在")
	}
	if article.AuthorID != userID && role != "admin" {
		return errors.New("无权限操作")
	}
	err = s.articleRepo.Delete(id)
	if err == nil {
		s.invalidateArticleCaches(id)
	}
	return err
}

func (s *ArticleService) Get(id uint) (*model.Article, error) {
	// 尝试从缓存读取
	cacheKey := cache.KeyArticleDetail(id)
	if cached, err := cache.Get(cacheKey); err == nil {
		var article model.Article
		if json.Unmarshal([]byte(cached), &article) == nil {
			// 缓存命中，异步增加浏览量
			go s.articleRepo.IncrementViewCount(id)
			return &article, nil
		}
	}

	// 缓存未命中
	_ = s.articleRepo.IncrementViewCount(id)
	article, err := s.articleRepo.FindByIDWithDetails(id)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if data, err := json.Marshal(article); err == nil {
		cache.Set(cacheKey, string(data), cache.TTLArticleDetail)
	}

	return article, nil
}

func (s *ArticleService) List(params repository.ArticleListParams) ([]model.Article, int64, error) {
	return s.articleRepo.List(params)
}

func (s *ArticleService) ListPopular(limit int) ([]model.Article, error) {
	// 尝试从缓存读取
	cacheKey := cache.KeyArticlesPopular(limit)
	if cached, err := cache.Get(cacheKey); err == nil {
		var articles []model.Article
		if json.Unmarshal([]byte(cached), &articles) == nil {
			return articles, nil
		}
	}

	articles, err := s.articleRepo.ListPopular(limit)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(articles); err == nil {
		cache.Set(cacheKey, string(data), cache.TTLArticlesPopular)
	}

	return articles, nil
}

func (s *ArticleService) ListRecent(limit int) ([]model.Article, error) {
	// 尝试从缓存读取
	cacheKey := cache.KeyArticlesRecent(limit)
	if cached, err := cache.Get(cacheKey); err == nil {
		var articles []model.Article
		if json.Unmarshal([]byte(cached), &articles) == nil {
			return articles, nil
		}
	}

	articles, err := s.articleRepo.ListRecent(limit)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(articles); err == nil {
		cache.Set(cacheKey, string(data), cache.TTLArticlesRecent)
	}

	return articles, nil
}

func (s *ArticleService) ListByUserID(userID uint) ([]model.Article, error) {
	return s.articleRepo.ListByUserID(userID)
}

func (s *ArticleService) ListPending() ([]model.Article, error) {
	return s.articleRepo.ListPending()
}

func (s *ArticleService) AdminDelete(id uint) error {
	err := s.articleRepo.Delete(id)
	if err == nil {
		s.invalidateArticleCaches(id)
	}
	return err
}

func (s *ArticleService) AdminUpdate(id uint, input UpdateArticleInput) (*model.Article, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	updates := map[string]interface{}{
		"title":   input.Title,
		"content": input.Content,
		"summary": input.Summary,
		"cover":   input.Cover,
	}
	if input.CategoryID != 0 {
		updates["category_id"] = input.CategoryID
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	if input.Status == "published" && article.PublishTime == nil {
		now := time.Now()
		updates["publish_time"] = &now
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.articleRepo.UpdateWithTx(tx, article, updates); err != nil {
			return err
		}
		if err := s.articleRepo.ReplaceTags(tx, article.ID, input.TagIDs); err != nil {
			return err
		}
		return nil
	})

	if err == nil {
		s.invalidateArticleCaches(id)
	}

	return article, err
}

func (s *ArticleService) AdminReview(id uint, status, reason string) error {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return errors.New("文章不存在")
	}
	if status != "published" && status != "rejected" {
		return errors.New("状态错误")
	}
	if status == "rejected" && reason == "" {
		return errors.New("驳回时必须提供原因")
	}

	updates := map[string]interface{}{
		"status": status,
	}
	if status == "published" && article.PublishTime == nil {
		now := time.Now()
		updates["publish_time"] = &now
	}
	if err := s.articleRepo.Update(article, updates); err != nil {
		return err
	}

	// 失效缓存
	s.invalidateArticleCaches(id)

	statusText := "通过"
	if status == "rejected" {
		statusText = "被驳回"
	}
	if err := s.notifRepo.Create(&model.Notification{
		UserID: article.AuthorID,
		Type:   "review",
		Title:  "文章审核" + statusText,
		Content: fmt.Sprintf("你的文章《%s》审核%s%s", article.Title, statusText, func() string {
			if status == "rejected" {
				return ", 原因: " + reason
			}
			return ""
		}()),
		Data: fmt.Sprintf(`{"article_id":%d,"status":"%s"}`, article.ID, status),
	}); err != nil {
		log.Printf("发送审核通知失败: %v", err)
	}

	return nil
}

func (s *ArticleService) AdminList(params repository.AdminArticleListParams) ([]model.Article, int64, error) {
	return s.articleRepo.AdminList(params)
}

// invalidateArticleCaches 失效单篇文章相关缓存
func (s *ArticleService) invalidateArticleCaches(id uint) {
	cache.Delete(cache.KeyArticleDetail(id))
	s.invalidateListCaches()
}

// invalidateListCaches 失效文章列表相关缓存
func (s *ArticleService) invalidateListCaches() {
	cache.DeleteByPattern(cache.KeyArticlesPopularPrefix + "*")
	cache.DeleteByPattern(cache.KeyArticlesRecentPrefix + "*")
}
