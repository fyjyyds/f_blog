package handler

import (
	"f_blog/backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListBanners(c *gin.Context) {
	banners, err := service.App.Banner.ListActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"banners": banners})
}
