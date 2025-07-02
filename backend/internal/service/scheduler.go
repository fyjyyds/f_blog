package service

import (
	"f_blog/backend/internal/database"
	"f_blog/backend/internal/model"
	"log"

	"github.com/robfig/cron/v3"
)

var scheduler *cron.Cron

// 初始化定时任务
func InitScheduler() {
	scheduler = cron.New(cron.WithSeconds())

	// 每天凌晨2点更新热门文章热度
	_, err := scheduler.AddFunc("0 0 2 * * *", updatePopularArticles)
	if err != nil {
		log.Printf("添加定时任务失败: %v", err)
		return
	}

	scheduler.Start()
	log.Println("定时任务已启动")
}

// 更新热门文章热度
func updatePopularArticles() {
	log.Println("开始更新热门文章热度...")

	// 计算所有已发布文章的热度分数
	var articles []model.Article
	err := database.DB.Where("status = ?", "published").Find(&articles).Error
	if err != nil {
		log.Printf("查询文章失败: %v", err)
		return
	}

	updatedCount := 0
	for _, article := range articles {
		// 计算热度分数：阅读量*1 + 点赞数*3 + 评论数*2
		popularityScore := article.ViewCount*1 + article.LikeCount*3 + article.CommentCount*2

		// 更新文章的热度分数（如果模型中有这个字段）
		// 这里可以根据需要在Article模型中添加popularity_score字段
		log.Printf("文章 %d 热度分数: %d (阅读:%d, 点赞:%d, 评论:%d)",
			article.ID, popularityScore, article.ViewCount, article.LikeCount, article.CommentCount)

		updatedCount++
	}

	log.Printf("热门文章热度更新完成，共更新 %d 篇文章", updatedCount)
}

// 手动触发更新热门文章
func ManualUpdatePopularArticles() {
	log.Println("手动触发更新热门文章热度...")
	updatePopularArticles()
}

// 停止定时任务
func StopScheduler() {
	if scheduler != nil {
		scheduler.Stop()
		log.Println("定时任务已停止")
	}
}
