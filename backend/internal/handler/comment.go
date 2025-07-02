package handler

import (
	"f_blog/backend/internal/database"
	"f_blog/backend/internal/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 发表评论
func CreateComment(c *gin.Context) {
	articleID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")
	var req struct {
		Content     string `json:"content" binding:"required"`
		ParentID    uint   `json:"parent_id"`
		ReplyToUser string `json:"reply_to_user"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comment := model.Comment{
		ArticleID:   uint(articleID),
		UserID:      userID,
		ParentID:    req.ParentID,
		Content:     req.Content,
		Status:      "approved",
		ReplyToUser: req.ReplyToUser,
	}
	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论失败"})
		return
	}
	// 评论成功后，文章评论数+1
	database.DB.Model(&model.Article{}).Where("id = ?", articleID).Update("comment_count", gorm.Expr("comment_count + 1"))

	// 评论成功后，给文章作者发送通知（不通知自己）
	var article model.Article
	database.DB.First(&article, articleID)
	if article.AuthorID != userID {
		database.DB.Create(&model.Notification{
			UserID:  article.AuthorID,
			Type:    "comment",
			Title:   "你的文章有新评论",
			Content: "有人评论了你的文章：" + article.Title,
			Data:    fmt.Sprintf(`{"article_id":%d}`, article.ID),
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "评论成功", "data": comment})
}

// 获取评论列表
func ListComments(c *gin.Context) {
	articleID, _ := strconv.Atoi(c.Param("id"))
	var comments []model.Comment
	if err := database.DB.
		Where("article_id = ?", articleID).
		Order("created_at asc").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	// 查询所有相关用户信息
	userIDs := make([]uint, 0, len(comments))
	userIDSet := make(map[uint]struct{})
	for _, c := range comments {
		if _, ok := userIDSet[c.UserID]; !ok {
			userIDs = append(userIDs, c.UserID)
			userIDSet[c.UserID] = struct{}{}
		}
	}
	var users []model.User
	database.DB.Where("id IN ?", userIDs).Find(&users)
	userMap := make(map[uint]model.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	// 构造带 display_name 的评论返回
	type CommentWithDisplayName struct {
		model.Comment
		DisplayName string `json:"display_name"`
	}
	result := make([]CommentWithDisplayName, 0, len(comments))
	for _, c := range comments {
		user := userMap[c.UserID]
		displayName := user.Nickname
		if displayName == "" {
			displayName = user.Username
		}
		result = append(result, CommentWithDisplayName{
			Comment:     c,
			DisplayName: displayName,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// 删除评论
func DeleteComment(c *gin.Context) {
	commentID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")
	userRole := c.GetString("role")
	var comment model.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}
	if comment.UserID != userID && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限操作"})
		return
	}
	if err := database.DB.Delete(&model.Comment{}, commentID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	// 删除评论后，文章评论数-1
	database.DB.Model(&model.Article{}).Where("id = ?", comment.ArticleID).Update("comment_count", gorm.Expr("comment_count - 1"))
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
