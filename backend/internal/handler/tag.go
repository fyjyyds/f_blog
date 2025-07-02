package handler

import (
	"f_blog/backend/internal/database"
	"f_blog/backend/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListTags(c *gin.Context) {
	var tags []model.Tag
	if err := database.DB.Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, tags)
}


