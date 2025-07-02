package handler

import (
	"f_blog/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.RouterGroup) {
	r.POST("/like", middleware.JWTAuth(), Like)
	r.DELETE("/like", middleware.JWTAuth(), Unlike)
	r.GET("/like/status", middleware.JWTAuth(), LikeStatus)

	userGroup := r.Group("/user")
	userGroup.Use(middleware.JWTAuth())
	{
		userGroup.POST("/change-password", ChangePassword)
	}
}
