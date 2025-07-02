package database

import (
	"f_blog/backend/internal/config"
	"fmt"
	"log"
	"strings"

	"f_blog/backend/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
	)

	// 配置GORM日志级别，减少不必要的日志输出
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error), // 只显示错误日志
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return err
	}
	log.Println("Database connected successfully")

	// 自动迁移所有模型，忽略约束删除错误
	if err := DB.AutoMigrate(
		&model.User{},
		&model.Article{},
		&model.Category{},
		&model.Tag{},
		&model.Comment{},
		&model.Like{},
		&model.Follow{},
		&model.Banner{},
		&model.ArticleTag{},
		&model.Setting{},
	); err != nil {
		// 检查是否是约束删除错误，如果是则忽略
		if strings.Contains(err.Error(), "Can't DROP") && strings.Contains(err.Error(), "check that column/key exists") {
			log.Printf("Warning: Ignoring constraint drop error: %v", err)
		} else {
			return fmt.Errorf("failed to initialize database: %v", err)
		}
	}

	return nil
}
