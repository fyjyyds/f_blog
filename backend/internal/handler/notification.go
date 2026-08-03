package handler

import (
	"f_blog/backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")
	notifications, err := service.App.Notification.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

func GetNotification(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")
	notification, err := service.App.Notification.Get(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在"})
		return
	}
	c.JSON(http.StatusOK, notification)
}

func MarkNotificationRead(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")
	if err := service.App.Notification.MarkRead(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已标记为已读"})
}

func MarkAllNotificationsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := service.App.Notification.MarkAllRead(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "全部已读"})
}

func DeleteNotification(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")
	if err := service.App.Notification.Delete(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}
