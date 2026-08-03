package handler

import (
	"f_blog/backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListCategories(c *gin.Context) {
	categories, err := service.App.Category.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, categories)
}
