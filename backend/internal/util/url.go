package util

import "strings"

// StaticPrefix 静态文件 URL 前缀
const StaticPrefix = "/static/"

// StaticURL 将相对路径拼接为完整的静态文件 URL
// 输入: "banners/banner1.jpg" → 输出: "/static/banners/banner1.jpg"
// 输入: "/static/banners/banner1.jpg" → 输出: "/static/banners/banner1.jpg" (已带前缀则不重复)
func StaticURL(path string) string {
	if path == "" {
		return ""
	}
	// 已经带了 /static/ 前缀，直接返回
	if strings.HasPrefix(path, StaticPrefix) || strings.HasPrefix(path, "/static/") {
		return path
	}
	// 去掉开头的 / 避免双斜杠
	path = strings.TrimPrefix(path, "/")
	return StaticPrefix + path
}

// RelativePath 从完整 URL 提取相对路径（存库用）
// 输入: "/static/banners/banner1.jpg" → 输出: "banners/banner1.jpg"
func RelativePath(url string) string {
	if strings.HasPrefix(url, StaticPrefix) {
		return strings.TrimPrefix(url, StaticPrefix)
	}
	if strings.HasPrefix(url, "/static/") {
		return strings.TrimPrefix(url, "/static/")
	}
	return url
}
