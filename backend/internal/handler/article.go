package handler

import (
	"f_blog/backend/internal/database"
	"f_blog/backend/internal/model"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 发布文章
func CreateArticle(c *gin.Context) {
	var req struct {
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content" binding:"required"`
		Summary    string `json:"summary"`
		Cover      string `json:"cover"`
		CategoryID uint   `json:"category_id"`
		Status     string `json:"status"`
		TagIDs     []uint `json:"tag_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRole := c.GetString("role")
	status := "pending"
	if userRole == "admin" {
		if req.Status != "" {
			status = req.Status
		} else {
			status = "published"
		}
	}

	authorID := c.GetUint("user_id")
	article := model.Article{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Cover:      req.Cover,
		AuthorID:   authorID,
		CategoryID: req.CategoryID,
		Status:     status,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if status == "published" {
		now := time.Now()
		article.PublishTime = &now
	}
	if err := database.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布失败"})
		return
	}
	for _, tagID := range req.TagIDs {
		database.DB.Create(&model.ArticleTag{
			ArticleID: article.ID,
			TagID:     tagID,
			CreatedAt: time.Now(),
		})
	}
	c.JSON(http.StatusOK, article)
}

// 编辑文章
func UpdateArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		Summary    string `json:"summary"`
		Cover      string `json:"cover"`
		CategoryID uint   `json:"category_id"`
		Status     string `json:"status"`
		TagIDs     []uint `json:"tag_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var article model.Article
	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	userID := c.GetUint("user_id")
	userRole := c.GetString("role")
	if article.AuthorID != userID && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限操作"})
		return
	}
	updates := map[string]interface{}{
		"title":       req.Title,
		"content":     req.Content,
		"summary":     req.Summary,
		"cover":       req.Cover,
		"category_id": req.CategoryID,
		"status":      req.Status,
		"updated_at":  time.Now(),
	}
	if req.Status == "published" && article.PublishTime == nil {
		now := time.Now()
		updates["publish_time"] = &now
	}
	if err := database.DB.Model(&article).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	database.DB.Where("article_id = ?", article.ID).Delete(&model.ArticleTag{})
	for _, tagID := range req.TagIDs {
		database.DB.Create(&model.ArticleTag{
			ArticleID: article.ID,
			TagID:     tagID,
			CreatedAt: time.Now(),
		})
	}
	c.JSON(http.StatusOK, article)
}

// 删除文章
func DeleteArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var article model.Article
	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	userID := c.GetUint("user_id")
	userRole := c.GetString("role")
	if article.AuthorID != userID && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限操作"})
		return
	}
	if err := database.DB.Delete(&model.Article{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 文章列表
func ListArticles(c *gin.Context) {
	var articles []model.Article

	// 1. 先构造基础筛选条件
	baseDB := database.DB.Where("status = ?", "published")

	if categoryID := c.Query("category_id"); categoryID != "" {
		baseDB = baseDB.Where("category_id = ?", categoryID)
	}
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	if startTime != "" && endTime != "" {
		baseDB = baseDB.Where("created_at BETWEEN ? AND ?", startTime, endTime)
	} else if startTime != "" {
		baseDB = baseDB.Where("created_at >= ?", startTime)
	} else if endTime != "" {
		baseDB = baseDB.Where("created_at <= ?", endTime)
	}

	// 2. 统计总数
	var total int64
	baseDB.Model(&model.Article{}).Count(&total)

	// 3. 再基于baseDB做分页、排序、预加载、Join等
	db := baseDB

	// 标签筛选（支持多个标签）
	if tagIDs := c.Query("tag_ids"); tagIDs != "" {
		ids := strings.Split(tagIDs, ",")
		db = db.Joins("JOIN article_tags at ON at.article_id = articles.id").
			Where("at.tag_id IN ?", ids)
	} else if tagID := c.Query("tag_id"); tagID != "" {
		db = db.Joins("JOIN article_tags at ON at.article_id = articles.id").
			Where("at.tag_id = ?", tagID)
	}

	// 排序
	sort := c.DefaultQuery("sort", "new")
	switch sort {
	case "hot":
		db = db.Order("view_count desc")
	case "comment":
		db = db.Order("comment_count desc")
	default:
		db = db.Order("created_at desc")
	}

	// 分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "8")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 8
	}

	// 分页
	db = db.Offset((page - 1) * pageSize).Limit(pageSize)

	// 预加载作者、分类、标签
	db = db.Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username", "nickname", "avatar")
	})
	db = db.Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	})
	db = db.Preload("Tags", func(db *gorm.DB) *gorm.DB {
		return db.Select("tags.id", "tags.name", "tags.color")
	})

	if err := db.Find(&articles).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"data":      articles,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// 获取热门文章
func GetPopularArticles(c *gin.Context) {
	limit := 10 // 默认返回10篇热门文章
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	var articles []model.Article
	db := database.DB

	// 只查询已发布的文章
	db = db.Where("status = ?", "published")

	// 按热度分数排序（MySQL 8.0+，不Select新字段）
	db = db.Order("view_count * 1 + like_count * 3 + comment_count * 2 DESC")

	// 限制返回数量
	db = db.Limit(limit)

	// 预加载作者信息
	db = db.Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username", "nickname", "avatar")
	})

	// 预加载分类信息
	db = db.Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	})

	// 预加载标签信息
	db = db.Preload("Tags", func(db *gorm.DB) *gorm.DB {
		return db.Select("tags.id", "tags.name", "tags.color")
	})

	if err := db.Find(&articles).Error; err != nil {
		c.JSON(500, gin.H{"error": "查询热门文章失败"})
		return
	}

	// 在Go里计算热度分数
	for i := range articles {
		articles[i].Summary = "热度分数：" + strconv.Itoa(int(articles[i].ViewCount*1+articles[i].LikeCount*3+articles[i].CommentCount*2))
	}

	c.JSON(200, gin.H{
		"data":       articles,
		"count":      len(articles),
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// 文章详情
func GetArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var article model.Article
	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	// 访问量自增
	database.DB.Model(&model.Article{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + 1"))
	// 查询作者信息
	var author model.User
	database.DB.Select("id", "username", "nickname", "avatar").First(&author, article.AuthorID)

	// 查询标签
	var tags []model.Tag
	database.DB.Joins("JOIN article_tags at ON at.tag_id = tags.id").
		Where("at.article_id = ?", article.ID).
		Find(&tags)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":           article.ID,
			"title":        article.Title,
			"content":      article.Content,
			"summary":      article.Summary,
			"cover":        article.Cover,
			"category_id":  article.CategoryID,
			"status":       article.Status,
			"created_at":   article.CreatedAt,
			"updated_at":   article.UpdatedAt,
			"publish_time": article.PublishTime,
			"author": gin.H{
				"id":       author.ID,
				"username": author.Username,
				"nickname": author.Nickname,
				"avatar":   author.Avatar,
			},
			"tags": tags,
		},
	})
}

// 管理员审核文章（通过/驳回）
func AdminReviewArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Status string `json:"status" binding:"required"` // "published"/"rejected"
		Reason string `json:"reason"`                    // 驳回原因
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Status != "published" && req.Status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "状态错误"})
		return
	}

	// 如果是驳回状态但没有提供原因，返回错误
	if req.Status == "rejected" && req.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "驳回时必须提供原因"})
		return
	}

	var article model.Article
	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	updates := map[string]interface{}{
		"status":     req.Status,
		"updated_at": time.Now(),
	}
	if req.Status == "published" && article.PublishTime == nil {
		now := time.Now()
		updates["publish_time"] = &now
	}
	if err := database.DB.Model(&article).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "审核失败"})
		return
	}

	// 审核后，给作者发送通知
	statusText := "通过"
	if req.Status == "rejected" {
		statusText = "被驳回"
	}
	database.DB.Create(&model.Notification{
		UserID: article.AuthorID,
		Type:   "review",
		Title:  "文章审核" + statusText,
		Content: fmt.Sprintf("你的文章《%s》审核%s%s", article.Title, statusText, func() string {
			if req.Status == "rejected" {
				return ", 原因: " + req.Reason
			} else {
				return ""
			}
		}()),
		Data: fmt.Sprintf(`{"article_id":%d,"status":"%s"}`, article.ID, req.Status),
	})

	c.JSON(http.StatusOK, gin.H{"message": "审核成功"})
}

// 管理员删除文章（不限制作者）
func AdminDeleteArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := database.DB.Delete(&model.Article{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 管理员编辑文章（不限制作者）
func AdminUpdateArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		Summary    string `json:"summary"`
		Cover      string `json:"cover"`
		CategoryID uint   `json:"category_id"`
		Status     string `json:"status"`
		TagIDs     []uint `json:"tag_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var article model.Article
	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	updates := map[string]interface{}{
		"title":       req.Title,
		"content":     req.Content,
		"summary":     req.Summary,
		"cover":       req.Cover,
		"category_id": req.CategoryID,
		"status":      req.Status,
		"updated_at":  time.Now(),
	}
	if req.Status == "published" && article.PublishTime == nil {
		now := time.Now()
		updates["publish_time"] = &now
	}
	if err := database.DB.Model(&article).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	database.DB.Where("article_id = ?", article.ID).Delete(&model.ArticleTag{})
	for _, tagID := range req.TagIDs {
		database.DB.Create(&model.ArticleTag{
			ArticleID: article.ID,
			TagID:     tagID,
			CreatedAt: time.Now(),
		})
	}
	c.JSON(http.StatusOK, article)
}

// 获取待审核文章列表
func ListPendingArticles(c *gin.Context) {
	var articles []model.Article
	if err := database.DB.Where("status = ?", "pending").Order("created_at desc").Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, articles)
}

// 管理员获取所有文章列表（支持分页、搜索、过滤）
func AdminListArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	categoryID := c.Query("category_id")
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	db := database.DB.Model(&model.Article{})

	// 状态过滤
	if status != "" {
		db = db.Where("status = ?", status)
	}

	// 分类过滤
	if categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}

	// 搜索过滤
	if search != "" {
		db = db.Where("title LIKE ? OR content LIKE ? OR summary LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	var total int64
	db.Count(&total)

	// 获取文章列表
	var articles []model.Article
	if err := db.Order("created_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	// 构建响应数据
	var responseArticles []gin.H
	for _, article := range articles {
		// 查询作者信息
		var author model.User
		database.DB.Select("id, username, nickname, avatar").First(&author, article.AuthorID)

		// 查询分类信息
		var category model.Category
		database.DB.Select("id, name").First(&category, article.CategoryID)

		// 查询标签信息
		var tags []model.Tag
		database.DB.Joins("JOIN article_tags at ON at.tag_id = tags.id").
			Where("at.article_id = ?", article.ID).
			Select("tags.id, tags.name, tags.color").
			Find(&tags)

		responseArticles = append(responseArticles, gin.H{
			"id":           article.ID,
			"title":        article.Title,
			"content":      article.Content,
			"summary":      article.Summary,
			"cover":        article.Cover,
			"category_id":  article.CategoryID,
			"status":       article.Status,
			"view_count":   article.ViewCount,
			"like_count":   article.LikeCount,
			"created_at":   article.CreatedAt,
			"updated_at":   article.UpdatedAt,
			"publish_time": article.PublishTime,
			"author":       author,
			"category":     category,
			"tags":         tags,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":    responseArticles,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (int(total) + pageSize - 1) / pageSize,
	})
}
