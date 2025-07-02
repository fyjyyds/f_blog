package config

import (
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Charset  string
}

type JWTConfig struct {
	Secret     string
	ExpireTime string
}

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
}

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string // 授权码
	From     string
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "mysql"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getIntEnv("DB_PORT", 3306),
			Username: getEnv("DB_USERNAME", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_DATABASE", "f_blog"),
			Charset:  getEnv("DB_CHARSET", "utf8mb4"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key"),
			ExpireTime: getEnv("JWT_EXPIRE_TIME", "24h"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func LoadEmailConfigQQ() *EmailConfig {
	return &EmailConfig{
		Host:     getEnv("EMAIL_QQ_HOST", "smtp.qq.com"),
		Port:     getIntEnv("EMAIL_QQ_PORT", 465),
		Username: getEnv("EMAIL_QQ_USERNAME", ""),
		Password: getEnv("EMAIL_QQ_PASSWORD", ""),
		From:     getEnv("EMAIL_QQ_FROM", ""),
	}
}
func LoadEmailConfig163() *EmailConfig {
	return &EmailConfig{
		Host:     getEnv("EMAIL_163_HOST", "smtp.163.com"),
		Port:     getIntEnv("EMAIL_163_PORT", 465),
		Username: getEnv("EMAIL_163_USERNAME", ""),
		Password: getEnv("EMAIL_163_PASSWORD", ""),
		From:     getEnv("EMAIL_163_FROM", ""),
	}
}
