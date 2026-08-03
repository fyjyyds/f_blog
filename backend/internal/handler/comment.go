package handler

import (
	"f_blog/backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

	comment, err := service.App.Comment.Create(uint(articleID), userID, req.Content, req.ParentID, req.ReplyToUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "评论成功", "data": comment})
}

func ListComments(c *gin.Context) {
	articleID, _ := strconv.Atoi(c.Param("id"))
	result, err := service.App.Comment.ListByArticleID(uint(articleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func DeleteComment(c *gin.Context) {
	commentID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	if err := service.App.Comment.Delete(uint(commentID), userID, role); err != nil {
		if err.Error() == "无权限操作" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
