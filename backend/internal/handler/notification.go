package handler

import (
	"f_blog/backend/internal/database"
	"f_blog/backend/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取当前用户通知列表
func ListNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")
	var notifications []model.Notification
	database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications)
	c.JSON(http.StatusOK, notifications)
}

// 获取通知详情
func GetNotification(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")
	var notification model.Notification
	if err := database.DB.Where("user_id = ?", userID).First(&notification, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在"})
		return
	}
	c.JSON(http.StatusOK, notification)
}

// 标记通知为已读
func MarkNotificationRead(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")
	if err := database.DB.Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已标记为已读"})
}

// 全部标记为已读
func MarkAllNotificationsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := database.DB.Model(&model.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "全部已读"})
}

// 删除通知
func DeleteNotification(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Notification{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}
