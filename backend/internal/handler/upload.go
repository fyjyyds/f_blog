package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 允许的扩展名
var allowedExt = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".zip": true, ".rar": true, ".7z": true,
	".txt": true, ".md": true, ".py": true, ".js": true, ".ts": true,
	".go": true, ".java": true, ".cpp": true, ".c": true,
	".json": true, ".xml": true, ".csv": true,
	// 新增常见办公文件类型
	".doc": true, ".docx": true, ".pdf": true,
	".ppt": true, ".pptx": true, ".xls": true, ".xlsx": true,
}

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExt[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型"})
		return
	}
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过10MB"})
		return
	}
	dir := "static/upload"
	_ = os.MkdirAll(dir, os.ModePerm)
	filename := "file_" + strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	savePath := filepath.Join(dir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}
	url := "/static/upload/" + filename
	c.JSON(http.StatusOK, gin.H{"url": url})
}
