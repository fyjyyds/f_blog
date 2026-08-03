package app

import (
	"f_blog/backend/internal/cache"
	"f_blog/backend/internal/config"
	"f_blog/backend/internal/database"
	"f_blog/backend/internal/middleware"
	"f_blog/backend/internal/model"
	"f_blog/backend/internal/repository"
	"f_blog/backend/internal/service"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// InitApp 初始化应用程序
func InitApp(cfg *config.Config) error {
	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		return err
	}

	// 初始化 Redis
	if err := cache.Init(&cfg.Redis); err != nil {
		log.Printf("Warning: Redis init failed: %v, caching disabled", err)
		// Redis 初始化失败不阻断启动，降级为无缓存模式
	}

	// 初始化管理员账号
	initAdmin()

	// 设置JWT配置
	middleware.SetJWTConfig(&cfg.JWT)
	service.SetJWTConfig(&cfg.JWT)

	// 初始化 Repository 层
	userRepo := repository.NewUserRepository(database.DB)
	articleRepo := repository.NewArticleRepository(database.DB)
	commentRepo := repository.NewCommentRepository(database.DB)
	likeRepo := repository.NewLikeRepository(database.DB)
	followRepo := repository.NewFollowRepository(database.DB)
	notifRepo := repository.NewNotificationRepository(database.DB)
	categoryRepo := repository.NewCategoryRepository(database.DB)
	tagRepo := repository.NewTagRepository(database.DB)
	bannerRepo := repository.NewBannerRepository(database.DB)
	settingRepo := repository.NewSettingRepository(database.DB)

	// 初始化 Service 屡
	userSvc := service.NewUserService(userRepo)
	articleSvc := service.NewArticleService(database.DB, articleRepo, notifRepo)
	commentSvc := service.NewCommentService(database.DB, commentRepo, articleRepo, notifRepo, userRepo)
	likeSvc := service.NewLikeService(likeRepo, notifRepo, articleRepo, commentRepo)
	followSvc := service.NewFollowService(followRepo, userRepo, notifRepo)
	notifSvc := service.NewNotificationService(notifRepo)
	categorySvc := service.NewCategoryService(categoryRepo)
	tagSvc := service.NewTagService(tagRepo)
	bannerSvc := service.NewBannerService(bannerRepo)
	settingSvc := service.NewSettingService(settingRepo)

	// 注册到全局容器
	service.InitContainer(&service.Container{
		User:         userSvc,
		Article:      articleSvc,
		Comment:      commentSvc,
		Like:         likeSvc,
		Follow:       followSvc,
		Notification: notifSvc,
		Category:     categorySvc,
		Tag:          tagSvc,
		Banner:       bannerSvc,
		Setting:      settingSvc,
	})

	// 初始化定时任务
	service.InitScheduler()

	return nil
}

// initAdmin 初始化默认管理员账号
func initAdmin() {
	var count int64
	database.DB.Model(&model.User{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := model.User{
			Username:      "admin",
			Email:         "admin@example.com",
			Password:      string(hash),
			Role:          "admin",
			Status:        "active",
			EmailVerified: true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		database.DB.Create(&admin)
		log.Println("默认管理员账号 admin/admin123 已创建")
	}
}
