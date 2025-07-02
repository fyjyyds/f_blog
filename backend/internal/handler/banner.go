package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListBanners(c *gin.Context) {
	banners := []string{
		"/static/banners/banner1.jpg",
		"/static/banners/banner2.jpg",
		"/static/banners/banner3.jpg",
	}
	c.JSON(http.StatusOK, gin.H{"banners": banners})
}
