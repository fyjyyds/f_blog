package handler

import (
	"f_blog/backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListTags(c *gin.Context) {
	tags, err := service.App.Tag.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, tags)
}
