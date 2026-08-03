package handler

import (
	"f_blog/backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Like(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   uint   `json:"target_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := service.App.Like.Like(userID, req.TargetType, req.TargetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "点赞成功", "liked": result.Liked, "count": result.Count})
}

func Unlike(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		TargetType string `json:"target_type" form:"target_type" binding:"required"`
		TargetID   uint   `json:"target_id" form:"target_id" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := service.App.Like.Unlike(userID, req.TargetType, req.TargetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消点赞失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "取消点赞成功", "liked": result.Liked, "count": result.Count})
}

func LikeStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	typeStr := c.Query("target_type")
	idStr := c.Query("target_id")
	if typeStr == "" || idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数缺失"})
		return
	}
	targetID, _ := strconv.Atoi(idStr)

	result, err := service.App.Like.Status(userID, typeStr, uint(targetID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"liked": result.Liked, "count": result.Count})
}
