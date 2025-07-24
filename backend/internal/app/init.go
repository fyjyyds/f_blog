package app

import (
	"f_blog/backend/internal/config"
	"f_blog/backend/internal/database"
	"f_blog/backend/internal/middleware"
	"f_blog/backend/internal/model"
	"f_blog/backend/internal/service"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// InitApp 初始化应用程序
func InitApp(cfg *config.Config) error {
	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		return err
	}

	// 初始化管理员账号
	initAdmin()

	// 设置JWT配置
	middleware.SetJWTConfig(&cfg.JWT)
	service.SetJWTConfig(&cfg.JWT)

	// 初始化定时任务
	service.InitScheduler()

	return nil
}

// initAdmin 初始化默认管理员账号
func initAdmin() {
	var count int64
	database.DB.Model(&model.User{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := model.User{
			Username:      "admin",
			Email:         "admin@example.com",
			Password:      string(hash),
			Role:          "admin",
			Status:        "active",
			EmailVerified: true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		database.DB.Create(&admin)
		log.Println("默认管理员账号 admin/admin123 已创建")
	}
}
