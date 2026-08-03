package handler

import (
	"f_blog/backend/internal/repository"
	"f_blog/backend/internal/service"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

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

	article, err := service.App.Article.Create(service.CreateArticleInput{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Cover:      req.Cover,
		CategoryID: req.CategoryID,
		Status:     req.Status,
		TagIDs:     req.TagIDs,
		AuthorID:   c.GetUint("user_id"),
		Role:       c.GetString("role"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布失败"})
		return
	}
	c.JSON(http.StatusOK, article)
}

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

	article, err := service.App.Article.Update(uint(id), service.UpdateArticleInput{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Cover:      req.Cover,
		CategoryID: req.CategoryID,
		Status:     req.Status,
		TagIDs:     req.TagIDs,
	}, c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		if err.Error() == "无权限操作" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, article)
}

func DeleteArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	if err := service.App.Article.Delete(uint(id), userID, role); err != nil {
		if err.Error() == "无权限操作" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func ListArticles(c *gin.Context) {
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

	var categoryID uint
	if catIDStr := c.Query("category_id"); catIDStr != "" {
		if id, err := strconv.ParseUint(catIDStr, 10, 32); err == nil {
			categoryID = uint(id)
		}
	}

	var tagIDs []uint
	if tagIDsStr := c.Query("tag_ids"); tagIDsStr != "" {
		for _, s := range strings.Split(tagIDsStr, ",") {
			if id, err := strconv.ParseUint(strings.TrimSpace(s), 10, 32); err == nil {
				tagIDs = append(tagIDs, uint(id))
			}
		}
	} else if tagID := c.Query("tag_id"); tagID != "" {
		if id, err := strconv.ParseUint(tagID, 10, 32); err == nil {
			tagIDs = []uint{uint(id)}
		}
	}

	articles, total, err := service.App.Article.List(repository.ArticleListParams{
		Page:       page,
		PageSize:   pageSize,
		CategoryID: categoryID,
		StartTime:  c.Query("start_time"),
		EndTime:    c.Query("end_time"),
		TagIDs:     tagIDs,
		Sort:       c.DefaultQuery("sort", "new"),
	})
	if err != nil {
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

func GetPopularArticles(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	articles, err := service.App.Article.ListPopular(limit)
	if err != nil {
		c.JSON(500, gin.H{"error": "查询热门文章失败"})
		return
	}

	c.JSON(200, gin.H{
		"data":       articles,
		"count":      len(articles),
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	})
}

func GetArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	article, err := service.App.Article.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": article})
}

func AdminReviewArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Status string `json:"status" binding:"required"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.App.Article.AdminReview(uint(id), req.Status, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "审核成功"})
}

func AdminDeleteArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.App.Article.AdminDelete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

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

	article, err := service.App.Article.AdminUpdate(uint(id), service.UpdateArticleInput{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Cover:      req.Cover,
		CategoryID: req.CategoryID,
		Status:     req.Status,
		TagIDs:     req.TagIDs,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, article)
}

func ListPendingArticles(c *gin.Context) {
	articles, err := service.App.Article.ListPending()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, articles)
}

func AdminListArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var categoryID uint
	if catIDStr := c.Query("category_id"); catIDStr != "" {
		if id, err := strconv.ParseUint(catIDStr, 10, 32); err == nil {
			categoryID = uint(id)
		}
	}

	articles, total, err := service.App.Article.AdminList(repository.AdminArticleListParams{
		Page:       page,
		PageSize:   pageSize,
		Status:     c.Query("status"),
		CategoryID: categoryID,
		Search:     c.Query("search"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":    articles,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (int(total) + pageSize - 1) / pageSize,
	})
}
