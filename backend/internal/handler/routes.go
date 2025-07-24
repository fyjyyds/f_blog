package handler

import (
	"f_blog/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine) {
	// 静态文件服务
	r.Static("/static", "./static")

	// API路由组
	api := r.Group("/api/v1")
	{
		// 公开接口
		setupPublicRoutes(api)

		// 用户相关接口
		setupUserRoutes(api)

		// 管理员接口
		setupAdminRoutes(api)

		// 文章相关接口
		setupArticleRoutes(api)

		// 点赞相关接口（原RegisterRouter内容）
		api.POST("/like", middleware.JWTAuth(), Like)
		api.DELETE("/like", middleware.JWTAuth(), Unlike)
		api.GET("/like/status", middleware.JWTAuth(), LikeStatus)

		// 用户改密接口（原RegisterRouter内容）
		userGroup := api.Group("/user")
		userGroup.Use(middleware.JWTAuth())
		{
			userGroup.POST("/change-password", ChangePassword)
		}
	}
}

// setupPublicRoutes 设置公开路由
func setupPublicRoutes(api *gin.RouterGroup) {
	api.GET("/captcha", GetCaptcha)
	api.POST("/auth/reg ister", Register)
	api.POST("/auth/login", Login)
	api.GET("/banners", ListBanners)
	api.GET("/categories", ListCategories)
	api.GET("/tags", ListTags)
	api.GET("/activate", ActivateEmail)

	// 评论相关接口
	api.POST("/articles/:id/comments", middleware.JWTAuth(), CreateComment)
	api.GET("/articles/:id/comments", ListComments)
	api.DELETE("/comments/:id", middleware.JWTAuth(), DeleteComment)

	// 文件上传接口
	api.POST("/upload", middleware.JWTAuth(), UploadFile)
}

// setupUserRoutes 设置用户相关路由
func setupUserRoutes(api *gin.RouterGroup) {
	userGroup := api.Group("/user")
	userGroup.Use(middleware.JWTAuth())
	{
		userGroup.GET("/profile", GetProfile)
		userGroup.PUT("/profile", UpdateProfile)
		userGroup.POST("/avatar", UploadAvatar)
		userGroup.DELETE("/avatar", DeleteAvatar)
		userGroup.GET("/articles", ListMyArticles)
		userGroup.GET("/comments", ListMyComments)
		userGroup.GET("/likes", ListMyLikes)
		userGroup.GET("/followings", ListMyFollowings)
		userGroup.GET("/followers", ListMyFollowers)
		userGroup.POST("/follow", Follow)
		userGroup.DELETE("/follow", Unfollow)

		// 通知相关接口
		userGroup.GET("/notifications", ListNotifications)
		userGroup.GET("/notifications/:id", GetNotification)
		userGroup.POST("/notifications/:id/read", MarkNotificationRead)
		userGroup.POST("/notifications/read-all", MarkAllNotificationsRead)
		userGroup.DELETE("/notifications/:id", DeleteNotification)
	}
}

// setupAdminRoutes 设置管理员路由
func setupAdminRoutes(api *gin.RouterGroup) {
	adminGroup := api.Group("/admin")
	adminGroup.Use(middleware.JWTAuth(), middleware.AdminOnly())
	{
		// 用户管理
		adminGroup.GET("/users", AdminListUsers)
		adminGroup.PUT("/users/:id", AdminUpdateUser)
		adminGroup.POST("/users/:id/ban", AdminBanUser)
		adminGroup.POST("/users/:id/unban", AdminUnbanUser)
		adminGroup.DELETE("/users/:id", AdminDeleteUser)

		// 分类管理
		adminGroup.GET("/categories", AdminListCategories)
		adminGroup.POST("/categories", AdminCreateCategory)
		adminGroup.PUT("/categories/:id", AdminUpdateCategory)
		adminGroup.DELETE("/categories/:id", AdminDeleteCategory)

		// 标签管理
		adminGroup.GET("/tags", AdminListTags)
		adminGroup.POST("/tags", AdminCreateTag)
		adminGroup.PUT("/tags/:id", AdminUpdateTag)
		adminGroup.DELETE("/tags/:id", AdminDeleteTag)

		// 评论管理
		adminGroup.GET("/comments", AdminListComments)
		adminGroup.POST("/comments/:id/approve", AdminApproveComment)
		adminGroup.POST("/comments/:id/reject", AdminRejectComment)
		adminGroup.PUT("/comments/:id", AdminUpdateComment)
		adminGroup.DELETE("/comments/:id", AdminDeleteComment)

		// 横幅管理
		adminGroup.GET("/banners", AdminListBanners)
		adminGroup.POST("/banners", AdminCreateBanner)
		adminGroup.PUT("/banners/:id", AdminUpdateBanner)
		adminGroup.DELETE("/banners/:id", AdminDeleteBanner)

		// 文章管理
		adminGroup.GET("/articles", AdminListArticles)
		adminGroup.GET("/articles/pending", ListPendingArticles)
		adminGroup.POST("/articles/:id/review", AdminReviewArticle)
		adminGroup.PUT("/articles/:id", AdminUpdateArticle)
		adminGroup.DELETE("/articles/:id", AdminDeleteArticle)

		// 统计和活动
		adminGroup.GET("/stats", AdminGetStats)
		adminGroup.GET("/activities", AdminGetActivities)

		// 系统设置
		adminGroup.GET("/settings", AdminGetSettings)
		adminGroup.PUT("/settings", AdminUpdateSettings)
		adminGroup.POST("/settings/reset", AdminResetSettings)
		adminGroup.POST("/settings/test-email", AdminTestEmail)

		// 热门文章管理
		adminGroup.POST("/articles/popular/update", AdminUpdatePopularArticles)
	}
}

// setupArticleRoutes 设置文章相关路由
func setupArticleRoutes(api *gin.RouterGroup) {
	articleGroup := api.Group("/articles")
	{
		articleGroup.GET("", ListArticles)
		articleGroup.GET("/popular", GetPopularArticles)
		articleGroup.GET(":id", GetArticle)
		articleGroup.Use(middleware.JWTAuth())
		articleGroup.POST("", CreateArticle)
		articleGroup.PUT(":id", UpdateArticle)
		articleGroup.DELETE(":id", DeleteArticle)
	}
}
