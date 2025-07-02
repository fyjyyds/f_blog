package handler

import (
	"f_blog/backend/internal/database"
	"f_blog/backend/internal/model"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 更新目标对象的点赞数
func updateTargetLikeCount(targetType string, targetID uint) error {
	var count int64
	if err := database.DB.Model(&model.Like{}).Where("target_type = ? AND target_id = ?", targetType, targetID).Count(&count).Error; err != nil {
		return err
	}

	switch targetType {
	case "article":
		return database.DB.Model(&model.Article{}).Where("id = ?", targetID).Update("like_count", count).Error
	case "comment":
		return database.DB.Model(&model.Comment{}).Where("id = ?", targetID).Update("like_count", count).Error
	}
	return nil
}

// 点赞
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
	log.Println("Like user_id:", userID, req.TargetType, req.TargetID)
	like := model.Like{
		UserID:     userID,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
	}
	// 唯一索引防止重复点赞
	if err := database.DB.Where(&like).FirstOrCreate(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
		return
	}
	// 统计最新数量并更新目标对象的点赞数
	var count int64
	database.DB.Model(&model.Like{}).Where("target_type = ? AND target_id = ?", req.TargetType, req.TargetID).Count(&count)

	// 更新目标对象的点赞数
	if err := updateTargetLikeCount(req.TargetType, req.TargetID); err != nil {
		log.Printf("更新点赞数失败: %v", err)
	}

	// 点赞成功后，给目标作者发送通知（不通知自己）
	if req.TargetType == "article" {
		var article model.Article
		if err := database.DB.First(&article, req.TargetID).Error; err == nil && article.AuthorID != userID {
			database.DB.Create(&model.Notification{
				UserID:  article.AuthorID,
				Type:    "like",
				Title:   "你的文章被点赞",
				Content: "有人点赞了你的文章：" + article.Title,
				Data:    fmt.Sprintf(`{"article_id":%d}`, article.ID),
			})
		}
	} else if req.TargetType == "comment" {
		var comment model.Comment
		if err := database.DB.First(&comment, req.TargetID).Error; err == nil && comment.UserID != userID {
			database.DB.Create(&model.Notification{
				UserID:  comment.UserID,
				Type:    "like",
				Title:   "你的评论被点赞",
				Content: "有人点赞了你的评论",
				Data:    fmt.Sprintf(`{"comment_id":%d}`, comment.ID),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "点赞成功", "liked": true, "count": count})
}

// 取消点赞
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
	log.Println("Unlike user_id:", userID, req.TargetType, req.TargetID)
	if err := database.DB.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, req.TargetType, req.TargetID).Delete(&model.Like{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消点赞失败"})
		return
	}
	var count int64
	database.DB.Model(&model.Like{}).Where("target_type = ? AND target_id = ?", req.TargetType, req.TargetID).Count(&count)

	// 更新目标对象的点赞数
	if err := updateTargetLikeCount(req.TargetType, req.TargetID); err != nil {
		log.Printf("更新点赞数失败: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消点赞成功", "liked": false, "count": count})
}

// 获取点赞状态和数量
func LikeStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	typeStr := c.Query("target_type")
	idStr := c.Query("target_id")
	if typeStr == "" || idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数缺失"})
		return
	}
	targetID, _ := strconv.Atoi(idStr)
	var count int64
	database.DB.Model(&model.Like{}).Where("target_type = ? AND target_id = ?", typeStr, targetID).Count(&count)
	var liked bool
	if userID != 0 {
		// 使用Count查询替代First查询，避免record not found错误
		var userLikeCount int64
		database.DB.Model(&model.Like{}).Where("user_id = ? AND target_type = ? AND target_id = ?", userID, typeStr, targetID).Count(&userLikeCount)
		liked = userLikeCount > 0
	}
	c.JSON(http.StatusOK, gin.H{"liked": liked, "count": count})
}
