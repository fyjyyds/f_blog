package main

import (
	"f_blog/backend/internal/app"
	"f_blog/backend/internal/config"
	"f_blog/backend/internal/handler"
	"f_blog/backend/internal/middleware"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// 加载配置
	cfg := config.Load()

	// 初始化应用程序
	if err := app.InitApp(cfg); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	r.Use(middleware.SetupCORS())

	// 设置路由
	handler.SetupRoutes(r)

	// 启动服务器
	r.Run()
}
