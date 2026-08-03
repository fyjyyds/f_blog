package cache

import (
	"fmt"
	"time"
)

// 缓存 Key 常量
const (
	// 分类列表（极少变化）
	KeyCategoriesList = "categories:list"

	// 标签列表（极少变化）
	KeyTagsList = "tags:list"

	// 活跃横幅
	KeyBannersActive = "banners:active"

	// 文章详情（按 ID）
	KeyArticleDetailPrefix = "article:detail:"

	// 热门文章列表
	KeyArticlesPopularPrefix = "articles:popular:"

	// 最新文章列表
	KeyArticlesRecentPrefix = "articles:recent:"

	// 用户资料
	KeyUserProfilePrefix = "user:profile:"
)

// TTL 常量
const (
	TTLCategories   = 30 * time.Minute
	TTLTags         = 30 * time.Minute
	TTLBanners      = 15 * time.Minute
	TTLArticleDetail = 5 * time.Minute
	TTLArticlesPopular = 10 * time.Minute
	TTLArticlesRecent  = 3 * time.Minute
	TTLUserProfile     = 10 * time.Minute
)

// Key 生成函数

func KeyArticleDetail(id uint) string {
	return KeyArticleDetailPrefix + uintToStr(id)
}

func KeyArticlesPopular(limit int) string {
	return KeyArticlesPopularPrefix + intToStr(limit)
}

func KeyArticlesRecent(limit int) string {
	return KeyArticlesRecentPrefix + intToStr(limit)
}

func KeyUserProfile(userID uint) string {
	return KeyUserProfilePrefix + uintToStr(userID)
}

// 辅助函数
func uintToStr(n uint) string {
	return fmt.Sprintf("%d", n)
}

func intToStr(n int) string {
	return fmt.Sprintf("%d", n)
}
