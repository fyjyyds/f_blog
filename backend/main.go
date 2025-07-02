package main

import (
	"f_blog/backend/internal/config"
	"f_blog/backend/internal/database"
	"f_blog/backend/internal/handler"
	"f_blog/backend/internal/middleware"
	"f_blog/backend/internal/service"
	"log"

	"f_blog/backend/internal/model"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

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

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}
	cfg := config.Load()
	if err := database.Init(&cfg.Database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	initAdmin()

	// 设置JWT配置
	middleware.SetJWTConfig(&cfg.JWT)
	service.SetJWTConfig(&cfg.JWT)

	// 初始化定时任务
	service.InitScheduler()

	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Static("/static", "./static")

	api := r.Group("/api/v1")
	{
		api.GET("/captcha", handler.GetCaptcha)
		api.POST("/auth/register", handler.Register)
		api.POST("/auth/login", handler.Login)
		api.GET("/banners", handler.ListBanners)
		api.GET("/categories", handler.ListCategories)
		api.GET("/tags", handler.ListTags)
		api.GET("/activate", handler.ActivateEmail)
		// 评论相关接口
		api.POST("/articles/:id/comments", middleware.JWTAuth(), handler.CreateComment)
		api.GET("/articles/:id/comments", handler.ListComments)
		api.DELETE("/comments/:id", middleware.JWTAuth(), handler.DeleteComment)

		// 需要登录的接口
		api.POST("/upload", middleware.JWTAuth(), handler.UploadFile)

		// 需要登录的用户接口
		userGroup := api.Group("/user")
		userGroup.Use(middleware.JWTAuth())
		{
			userGroup.GET("/profile", handler.GetProfile)
			userGroup.PUT("/profile", handler.UpdateProfile)
			userGroup.POST("/avatar", handler.UploadAvatar)
			userGroup.DELETE("/avatar", handler.DeleteAvatar)
			userGroup.GET("/articles", handler.ListMyArticles)
			userGroup.GET("/comments", handler.ListMyComments)
			userGroup.GET("/likes", handler.ListMyLikes)
			userGroup.GET("/followings", handler.ListMyFollowings)
			userGroup.GET("/followers", handler.ListMyFollowers)
			userGroup.POST("/follow", handler.Follow)
			userGroup.DELETE("/follow", handler.Unfollow)

			// 通知相关接口
			userGroup.GET("/notifications", handler.ListNotifications)
			userGroup.GET("/notifications/:id", handler.GetNotification)
			userGroup.POST("/notifications/:id/read", handler.MarkNotificationRead)
			userGroup.POST("/notifications/read-all", handler.MarkAllNotificationsRead)
			userGroup.DELETE("/notifications/:id", handler.DeleteNotification)
		}

		// 仅管理员可访问的接口
		adminGroup := api.Group("/admin")
		adminGroup.Use(middleware.JWTAuth(), middleware.AdminOnly())
		{
			// 用户管理
			adminGroup.GET("/users", handler.AdminListUsers)
			adminGroup.PUT("/users/:id", handler.AdminUpdateUser)
			adminGroup.POST("/users/:id/ban", handler.AdminBanUser)
			adminGroup.POST("/users/:id/unban", handler.AdminUnbanUser)
			adminGroup.DELETE("/users/:id", handler.AdminDeleteUser)

			// 分类管理
			adminGroup.GET("/categories", handler.AdminListCategories)
			adminGroup.POST("/categories", handler.AdminCreateCategory)
			adminGroup.PUT("/categories/:id", handler.AdminUpdateCategory)
			adminGroup.DELETE("/categories/:id", handler.AdminDeleteCategory)

			// 标签管理
			adminGroup.GET("/tags", handler.AdminListTags)
			adminGroup.POST("/tags", handler.AdminCreateTag)
			adminGroup.PUT("/tags/:id", handler.AdminUpdateTag)
			adminGroup.DELETE("/tags/:id", handler.AdminDeleteTag)

			// 评论管理
			adminGroup.GET("/comments", handler.AdminListComments)
			adminGroup.POST("/comments/:id/approve", handler.AdminApproveComment)
			adminGroup.POST("/comments/:id/reject", handler.AdminRejectComment)
			adminGroup.PUT("/comments/:id", handler.AdminUpdateComment)
			adminGroup.DELETE("/comments/:id", handler.AdminDeleteComment)

			// 横幅管理
			adminGroup.GET("/banners", handler.AdminListBanners)
			adminGroup.POST("/banners", handler.AdminCreateBanner)
			adminGroup.PUT("/banners/:id", handler.AdminUpdateBanner)
			adminGroup.DELETE("/banners/:id", handler.AdminDeleteBanner)

			// 文章管理
			adminGroup.GET("/articles", handler.AdminListArticles)
			adminGroup.GET("/articles/pending", handler.ListPendingArticles)
			adminGroup.POST("/articles/:id/review", handler.AdminReviewArticle)
			adminGroup.PUT("/articles/:id", handler.AdminUpdateArticle)
			adminGroup.DELETE("/articles/:id", handler.AdminDeleteArticle)

			// 统计和活动
			adminGroup.GET("/stats", handler.AdminGetStats)
			adminGroup.GET("/activities", handler.AdminGetActivities)

			// 系统设置
			adminGroup.GET("/settings", handler.AdminGetSettings)
			adminGroup.PUT("/settings", handler.AdminUpdateSettings)
			adminGroup.POST("/settings/reset", handler.AdminResetSettings)
			adminGroup.POST("/settings/test-email", handler.AdminTestEmail)

			// 热门文章管理
			adminGroup.POST("/articles/popular/update", handler.AdminUpdatePopularArticles)
		}

		articleGroup := api.Group("/articles")
		{
			articleGroup.GET("", handler.ListArticles)
			articleGroup.GET("/popular", handler.GetPopularArticles)
			articleGroup.GET(":id", handler.GetArticle)
			articleGroup.Use(middleware.JWTAuth())
			articleGroup.POST("", handler.CreateArticle)
			articleGroup.PUT(":id", handler.UpdateArticle)
			articleGroup.DELETE(":id", handler.DeleteArticle)
		}

		handler.RegisterRouter(api)
	}
	r.Run()
}
