package handler

import (
	"f_blog/backend/internal/config"
	"f_blog/backend/internal/database"
	"f_blog/backend/internal/model"
	"f_blog/backend/internal/service"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 管理员获取用户列表
func AdminListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	db := database.DB.Model(&model.User{})

	if search != "" {
		db = db.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if status != "" {
		db = db.Where("status = ?", status)
	}

	var total int64
	db.Count(&total)

	var users []model.User
	if err := db.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&users).Error; err != nil {
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

// 管理员更新用户信息
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

	var user model.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	updates := map[string]interface{}{
		"username":   req.Username,
		"email":      req.Email,
		"nickname":   req.Nickname,
		"role":       req.Role,
		"status":     req.Status,
		"updated_at": time.Now(),
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// 管理员封禁用户
func AdminBanUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var user model.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	if err := database.DB.Model(&user).Update("status", "banned").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "封禁失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户已封禁"})
}

// 管理员解封用户
func AdminUnbanUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var user model.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	if err := database.DB.Model(&user).Update("status", "active").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解封失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户已解封"})
}

// 管理员删除用户
func AdminDeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := database.DB.Delete(&model.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户已删除"})
}

// 管理员获取分类列表
func AdminListCategories(c *gin.Context) {
	var categories []model.Category
	if err := database.DB.Order("sort_order asc, created_at desc").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// 管理员创建分类
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

	if req.Status == "" {
		req.Status = "active"
	}

	category := model.Category{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// 管理员更新分类
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

	var category model.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"sort_order":  req.SortOrder,
		"status":      req.Status,
		"updated_at":  time.Now(),
	}

	if err := database.DB.Model(&category).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// 管理员删除分类
func AdminDeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := database.DB.Delete(&model.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "分类已删除"})
}

// 管理员获取标签列表
func AdminListTags(c *gin.Context) {
	var tags []model.Tag
	if err := database.DB.Order("created_at desc").Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// 管理员创建标签
func AdminCreateTag(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Color string `json:"color"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Color == "" {
		req.Color = "#667eea"
	}

	tag := model.Tag{
		Name:      req.Name,
		Color:     req.Color,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// 管理员更新标签
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

	var tag model.Tag
	if err := database.DB.First(&tag, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}

	updates := map[string]interface{}{
		"name":       req.Name,
		"color":      req.Color,
		"updated_at": time.Now(),
	}

	if err := database.DB.Model(&tag).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// 管理员删除标签
func AdminDeleteTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := database.DB.Delete(&model.Tag{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "标签已删除"})
}

// 管理员获取评论列表
func AdminListComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	db := database.DB.Model(&model.Comment{})

	if status != "" {
		db = db.Where("status = ?", status)
	}

	var total int64
	db.Count(&total)

	var comments []model.Comment
	if err := db.Preload("User").Preload("Article").
		Offset(offset).Limit(pageSize).Order("created_at desc").Find(&comments).Error; err != nil {
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

// 管理员审核评论（通过）
func AdminApproveComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var comment model.Comment
	if err := database.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	if err := database.DB.Model(&comment).Update("status", "approved").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "审核失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评论已通过"})
}

// 管理员审核评论（驳回）
func AdminRejectComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var comment model.Comment
	if err := database.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	if err := database.DB.Model(&comment).Update("status", "rejected").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "审核失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评论已驳回"})
}

// 管理员更新评论
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

	var comment model.Comment
	if err := database.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	updates := map[string]interface{}{
		"content":    req.Content,
		"status":     req.Status,
		"updated_at": time.Now(),
	}

	if err := database.DB.Model(&comment).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// 管理员删除评论
func AdminDeleteComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := database.DB.Delete(&model.Comment{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评论已删除"})
}

// 管理员获取横幅列表
func AdminListBanners(c *gin.Context) {
	var banners []model.Banner
	if err := database.DB.Order("sort_order asc, created_at desc").Find(&banners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, banners)
}

// 管理员创建横幅
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

	if req.Status == "" {
		req.Status = "active"
	}

	banner := model.Banner{
		Title:     req.Title,
		Image:     req.Image,
		Link:      req.Link,
		SortOrder: req.SortOrder,
		Status:    req.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&banner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, banner)
}

// 管理员更新横幅
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
	var banner model.Banner
	if err := database.DB.First(&banner, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "横幅不存在"})
		return
	}
	updates := map[string]interface{}{
		"title":      req.Title,
		"image":      req.Image,
		"link":       req.Link,
		"sort_order": req.SortOrder,
		"status":     req.Status,
		"updated_at": time.Now(),
	}
	if err := database.DB.Model(&banner).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, banner)
}

// 管理员删除横幅
func AdminDeleteBanner(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := database.DB.Delete(&model.Banner{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "横幅已删除"})
}

// 管理员获取统计数据
func AdminGetStats(c *gin.Context) {
	var userCount, articleCount, commentCount, pendingArticleCount int64
	database.DB.Model(&model.User{}).Count(&userCount)
	database.DB.Model(&model.Article{}).Count(&articleCount)
	database.DB.Model(&model.Comment{}).Count(&commentCount)
	database.DB.Model(&model.Article{}).Where("status = ?", "pending").Count(&pendingArticleCount)
	c.JSON(http.StatusOK, gin.H{
		"user_count":            userCount,
		"article_count":         articleCount,
		"comment_count":         commentCount,
		"pending_article_count": pendingArticleCount,
	})
}

// 管理员获取最近活动
func AdminGetActivities(c *gin.Context) {
	limit := 10
	var recentArticles []model.Article
	database.DB.Order("created_at desc").Limit(limit).Find(&recentArticles)
	var recentComments []model.Comment
	database.DB.Order("created_at desc").Limit(limit).Find(&recentComments)
	c.JSON(http.StatusOK, gin.H{
		"recent_articles": recentArticles,
		"recent_comments": recentComments,
	})
}

// 管理员获取设置
func AdminGetSettings(c *gin.Context) {
	var settings []model.Setting
	if err := database.DB.Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取设置失败"})
		return
	}

	// 将设置转换为前端需要的格式
	basic := make(map[string]interface{})
	email := make(map[string]interface{})
	security := make(map[string]interface{})
	content := make(map[string]interface{})

	for _, setting := range settings {
		if strings.HasPrefix(setting.Key, "email_") {
			key := strings.TrimPrefix(setting.Key, "email_")
			email[key] = setting.Value
		} else if strings.HasPrefix(setting.Key, "security_") {
			key := strings.TrimPrefix(setting.Key, "security_")
			security[key] = setting.Value
		} else if strings.HasPrefix(setting.Key, "content_") {
			key := strings.TrimPrefix(setting.Key, "content_")
			content[key] = setting.Value
		} else {
			basic[setting.Key] = setting.Value
		}
	}

	// 如果没有设置，返回默认值
	if len(basic) == 0 {
		basic = map[string]interface{}{
			"siteName":        "F_Blog",
			"siteDescription": "一个现代化的博客系统",
			"siteKeywords":    "",
			"siteLogo":        "",
			"icp":             "",
		}
	}

	if len(email) == 0 {
		email = map[string]interface{}{
			"smtpHost":      "smtp.example.com",
			"smtpPort":      "587",
			"emailUser":     "admin@example.com",
			"emailPassword": "",
			"senderName":    "F_Blog",
		}
	}

	if len(security) == 0 {
		security = map[string]interface{}{
			"minPasswordLength": "8",
			"maxLoginAttempts":  "5",
			"lockoutDuration":   "30",
			"jwtExpireHours":    "24",
			"enableCaptcha":     "false",
			"enableTwoFactor":   "false",
		}
	}

	if len(content) == 0 {
		content = map[string]interface{}{
			"articlesPerPage":        "10",
			"summaryLength":          "200",
			"commentModeration":      "manual",
			"allowAnonymousComments": "false",
			"enableCommentLikes":     "true",
			"enableArticleLikes":     "true",
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"basic":    basic,
		"email":    email,
		"security": security,
		"content":  content,
	})
}

// 管理员更新设置
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

	// 保存基本设置
	if req.Basic != nil {
		for key, value := range req.Basic {
			if err := saveSetting(key, value, "string"); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "保存基本设置失败"})
				return
			}
		}
	}

	// 保存邮件设置
	if req.Email != nil {
		for key, value := range req.Email {
			if err := saveSetting("email_"+key, value, "string"); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "保存邮件设置失败"})
				return
			}
		}
	}

	// 保存安全设置
	if req.Security != nil {
		for key, value := range req.Security {
			if err := saveSetting("security_"+key, value, "string"); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "保存安全设置失败"})
				return
			}
		}
	}

	// 保存内容设置
	if req.Content != nil {
		for key, value := range req.Content {
			if err := saveSetting("content_"+key, value, "string"); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "保存内容设置失败"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "设置已更新"})
}

// 保存单个设置项
func saveSetting(key string, value interface{}, valueType string) error {
	log.Printf("saveSetting: key=%s, value=%v, type=%s", key, value, valueType)
	var setting model.Setting
	// 修正：key字段加反引号
	result := database.DB.Where("`key` = ?", key).First(&setting)
	if result.Error != nil {
		// 设置不存在，创建新的
		setting = model.Setting{
			Key:   key,
			Type:  valueType,
			Value: fmt.Sprintf("%v", value),
		}
		return database.DB.Create(&setting).Error
	} else {
		// 设置已存在，更新
		setting.Value = fmt.Sprintf("%v", value)
		setting.UpdatedAt = time.Now()
		return database.DB.Save(&setting).Error
	}
}

// 管理员重置设置
func AdminResetSettings(c *gin.Context) {
	// 清空所有设置
	if err := database.DB.Where("1 = 1").Delete(&model.Setting{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "重置设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "设置已重置为默认值"})
}

// 管理员测试邮件
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

	// 构建邮件配置
	emailCfg := &config.EmailConfig{
		Host:     req.SmtpHost,
		Port:     getIntFromString(req.SmtpPort, 587),
		Username: req.EmailUser,
		Password: req.EmailPassword,
		From:     req.SenderName,
	}

	// 发送测试邮件
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

// 辅助函数：将字符串转换为整数
func getIntFromString(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return defaultValue
}

// 手动更新热门文章热度
func AdminUpdatePopularArticles(c *gin.Context) {
	// 调用服务层的手动更新函数
	service.ManualUpdatePopularArticles()

	c.JSON(http.StatusOK, gin.H{
		"message":    "热门文章热度更新成功",
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	})
}
