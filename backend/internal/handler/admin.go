package handler

import (
	"f_blog/backend/internal/config"
	"f_blog/backend/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ========== 用户管理 ==========

func AdminListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	users, total, err := service.App.User.AdminList(page, pageSize, c.Query("search"), c.Query("status"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":       users,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (int(total) + pageSize - 1) / pageSize,
	})
}

func AdminUpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Nickname string `json:"nickname"`
		Role     string `json:"role"`
		Status   string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := service.App.User.AdminUpdate(uint(id), req.Username, req.Email, req.Nickname, req.Role, req.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func AdminBanUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.App.User.AdminBan(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "用户已封禁"})
}

func AdminUnbanUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.App.User.AdminUnban(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "用户已解封"})
}

func AdminDeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.App.User.AdminDelete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "用户已删除"})
}

// ========== 分类管理 ==========

func AdminListCategories(c *gin.Context) {
	categories, err := service.App.Category.AdminList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func AdminCreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
		Status      string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := service.App.Category.Create(req.Name, req.Description, req.SortOrder, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, category)
}

func AdminUpdateCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
		Status      string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := service.App.Category.Update(uint(id), req.Name, req.Description, req.SortOrder, req.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

func AdminDeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.App.Category.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "分类已删除"})
}

// ========== 标签管理 ==========

func AdminListTags(c *gin.Context) {
	tags, err := service.App.Tag.AdminList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, tags)
}

func AdminCreateTag(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Color string `json:"color"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := service.App.Tag.Create(req.Name, req.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, tag)
}

func AdminUpdateTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := service.App.Tag.Update(uint(id), req.Name, req.Color)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tag)
}

func AdminDeleteTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.App.Tag.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "标签已删除"})
}

// ========== 评论管理 ==========

func AdminListComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	comments, total, err := service.App.Comment.AdminList(page, pageSize, c.Query("status"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments":    comments,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (int(total) + pageSize - 1) / pageSize,
	})
}

func AdminApproveComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.App.Comment.AdminApprove(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "评论已通过"})
}

func AdminRejectComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.App.Comment.AdminReject(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "评论已驳回"})
}

func AdminUpdateComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Content string `json:"content"`
		Status  string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := service.App.Comment.AdminUpdate(uint(id), req.Content, req.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comment)
}

func AdminDeleteComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.App.Comment.AdminDelete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "评论已删除"})
}

// ========== 横幅管理 ==========

func AdminListBanners(c *gin.Context) {
	banners, err := service.App.Banner.AdminList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, banners)
}

func AdminCreateBanner(c *gin.Context) {
	var req struct {
		Title     string `json:"title" binding:"required"`
		Image     string `json:"image" binding:"required"`
		Link      string `json:"link"`
		SortOrder int    `json:"sort_order"`
		Status    string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	banner, err := service.App.Banner.Create(req.Title, req.Image, req.Link, req.SortOrder, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, banner)
}

func AdminUpdateBanner(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Title     string `json:"title"`
		Image     string `json:"image"`
		Link      string `json:"link"`
		SortOrder int    `json:"sort_order"`
		Status    string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	banner, err := service.App.Banner.Update(uint(id), req.Title, req.Image, req.Link, req.SortOrder, req.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, banner)
}

func AdminDeleteBanner(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.App.Banner.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "横幅已删除"})
}

// ========== 统计和活动 ==========

func AdminGetStats(c *gin.Context) {
	userCount, _ := service.App.User.Count()
	articleCount, _ := service.App.Article.Count()
	commentCount, _ := service.App.Comment.Count()
	pendingCount, _ := service.App.Article.CountByStatus("pending")

	c.JSON(http.StatusOK, gin.H{
		"user_count":            userCount,
		"article_count":         articleCount,
		"comment_count":         commentCount,
		"pending_article_count": pendingCount,
	})
}

func AdminGetActivities(c *gin.Context) {
	limit := 10
	recentArticles, _ := service.App.Article.ListRecent(limit)
	recentComments := []interface{}{}

	c.JSON(http.StatusOK, gin.H{
		"recent_articles": recentArticles,
		"recent_comments": recentComments,
	})
}

// ========== 系统设置 ==========

func AdminGetSettings(c *gin.Context) {
	settings, err := service.App.Setting.GetGrouped()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取设置失败"})
		return
	}
	c.JSON(http.StatusOK, settings)
}

func AdminUpdateSettings(c *gin.Context) {
	var req struct {
		Basic    map[string]interface{} `json:"basic"`
		Email    map[string]interface{} `json:"email"`
		Security map[string]interface{} `json:"security"`
		Content  map[string]interface{} `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.App.Setting.Update(req.Basic, req.Email, req.Security, req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存设置失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "设置已更新"})
}

func AdminResetSettings(c *gin.Context) {
	if err := service.App.Setting.Reset(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "重置设置失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "设置已重置为默认值"})
}

func AdminTestEmail(c *gin.Context) {
	var req struct {
		SmtpHost      string `json:"smtpHost"`
		SmtpPort      string `json:"smtpPort"`
		EmailUser     string `json:"emailUser"`
		EmailPassword string `json:"emailPassword"`
		SenderName    string `json:"senderName"`
		To            string `json:"to" binding:"required"`
		Subject       string `json:"subject"`
		Content       string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	from := req.SenderName
	if req.EmailUser != "" {
		from = req.SenderName + " <" + req.EmailUser + ">"
	}
	emailCfg := &config.EmailConfig{
		Host:     req.SmtpHost,
		Port:     getIntFromString(req.SmtpPort, 587),
		Username: req.EmailUser,
		Password: req.EmailPassword,
		From:     from,
	}

	subject := req.Subject
	if subject == "" {
		subject = "F_Blog 邮件测试"
	}
	content := req.Content
	if content == "" {
		content = "这是一封来自F_Blog系统的测试邮件，如果您收到这封邮件，说明邮件配置正确。"
	}

	if err := service.SendEmail(emailCfg, req.To, subject, content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "邮件发送失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "测试邮件发送成功"})
}

func AdminUpdatePopularArticles(c *gin.Context) {
	service.ManualUpdatePopularArticles()
	c.JSON(http.StatusOK, gin.H{
		"message":    "热门文章热度更新成功",
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	})
}
