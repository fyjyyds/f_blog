package repository

import (
	"f_blog/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// Count 统计文章总数
func (r *ArticleRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Article{}).Count(&count).Error
	return count, err
}

// CountByStatus 按状态统计文章数
func (r *ArticleRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Article{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func (r *ArticleRepository) FindByID(id uint) (*model.Article, error) {
	var article model.Article
	err := r.db.First(&article, id).Error
	return &article, err
}

func (r *ArticleRepository) FindByIDWithDetails(id uint) (*model.Article, error) {
	var article model.Article
	err := r.db.Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username", "nickname", "avatar")
	}).Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Preload("Tags", func(db *gorm.DB) *gorm.DB {
		return db.Select("tags.id", "tags.name", "tags.color")
	}).First(&article, id).Error
	return &article, err
}

func (r *ArticleRepository) Create(tx *gorm.DB, article *model.Article) error {
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	return tx.Create(article).Error
}

func (r *ArticleRepository) Update(article *model.Article, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return r.db.Model(article).Updates(updates).Error
}

func (r *ArticleRepository) UpdateWithTx(tx *gorm.DB, article *model.Article, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return tx.Model(article).Updates(updates).Error
}

func (r *ArticleRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 级联删除关联数据
		if err := tx.Where("article_id = ?", id).Delete(&model.ArticleTag{}).Error; err != nil {
			return err
		}
		if err := tx.Where("article_id = ?", id).Delete(&model.Comment{}).Error; err != nil {
			return err
		}
		if err := tx.Where("target_type = ? AND target_id = ?", "article", id).Delete(&model.Like{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Article{}, id).Error
	})
}

func (r *ArticleRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *ArticleRepository) List(params ArticleListParams) ([]model.Article, int64, error) {
	db := r.db.Model(&model.Article{}).Where("status = ?", "published")

	if params.CategoryID != 0 {
		db = db.Where("category_id = ?", params.CategoryID)
	}
	if params.StartTime != "" && params.EndTime != "" {
		db = db.Where("created_at BETWEEN ? AND ?", params.StartTime, params.EndTime)
	} else if params.StartTime != "" {
		db = db.Where("created_at >= ?", params.StartTime)
	} else if params.EndTime != "" {
		db = db.Where("created_at <= ?", params.EndTime)
	}
	if len(params.TagIDs) > 0 {
		db = db.Joins("JOIN article_tags at ON at.article_id = articles.id").
			Where("at.tag_id IN ?", params.TagIDs).
			Distinct("articles.id")
	}

	var total int64
	db.Count(&total)

	switch params.Sort {
	case "hot":
		db = db.Order("view_count desc")
	case "comment":
		db = db.Order("comment_count desc")
	default:
		db = db.Order("created_at desc")
	}

	db = db.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)
	db = db.Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username", "nickname", "avatar")
	})
	db = db.Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	})
	db = db.Preload("Tags", func(db *gorm.DB) *gorm.DB {
		return db.Select("tags.id", "tags.name", "tags.color")
	})

	var articles []model.Article
	err := db.Find(&articles).Error
	return articles, total, err
}

func (r *ArticleRepository) ListRecent(limit int) ([]model.Article, error) {
	var articles []model.Article
	err := r.db.Order("created_at desc").Limit(limit).Find(&articles).Error
	return articles, err
}

func (r *ArticleRepository) ListPopular(limit int) ([]model.Article, error) {
	var articles []model.Article
	err := r.db.Where("status = ?", "published").
		Order("view_count * 1 + like_count * 3 + comment_count * 2 DESC").
		Limit(limit).
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username", "nickname", "avatar")
		}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Select("tags.id", "tags.name", "tags.color")
		}).
		Find(&articles).Error
	return articles, err
}

func (r *ArticleRepository) ListByUserID(userID uint) ([]model.Article, error) {
	var articles []model.Article
	err := r.db.Where("author_id = ?", userID).Order("created_at desc").Find(&articles).Error
	return articles, err
}

func (r *ArticleRepository) ListPending() ([]model.Article, error) {
	var articles []model.Article
	err := r.db.Where("status = ?", "pending").Order("created_at desc").Find(&articles).Error
	return articles, err
}

// Tag operations on articles

func (r *ArticleRepository) ReplaceTags(tx *gorm.DB, articleID uint, tagIDs []uint) error {
	if err := tx.Where("article_id = ?", articleID).Delete(&model.ArticleTag{}).Error; err != nil {
		return err
	}
	for _, tagID := range tagIDs {
		if err := tx.Create(&model.ArticleTag{
			ArticleID: articleID,
			TagID:     tagID,
			CreatedAt: time.Now(),
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *ArticleRepository) FindTagsByArticleID(articleID uint) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Joins("JOIN article_tags at ON at.tag_id = tags.id").
		Where("at.article_id = ?", articleID).
		Find(&tags).Error
	return tags, err
}

// Admin operations

func (r *ArticleRepository) AdminList(params AdminArticleListParams) ([]model.Article, int64, error) {
	db := r.db.Model(&model.Article{})

	if params.Status != "" {
		db = db.Where("status = ?", params.Status)
	}
	if params.CategoryID != 0 {
		db = db.Where("category_id = ?", params.CategoryID)
	}
	if params.Search != "" {
		db = db.Where("title LIKE ? OR content LIKE ? OR summary LIKE ?",
			"%"+params.Search+"%", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	var total int64
	db.Count(&total)

	var articles []model.Article
	err := db.Order("created_at desc").
		Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username", "nickname", "avatar")
		}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Select("tags.id", "tags.name", "tags.color")
		}).
		Find(&articles).Error
	return articles, total, err
}

type ArticleListParams struct {
	Page       int
	PageSize   int
	CategoryID uint
	StartTime  string
	EndTime    string
	TagIDs     []uint
	Sort       string
}

type AdminArticleListParams struct {
	Page       int
	PageSize   int
	Status     string
	CategoryID uint
	Search     string
}
