package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"f_blog/backend/internal/config"
	"f_blog/backend/internal/database"
	"f_blog/backend/internal/model"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtConfig *config.JWTConfig

func SetJWTConfig(cfg *config.JWTConfig) {
	jwtConfig = cfg
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func generateToken(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func RegisterUser(db *gorm.DB, username, email, password string, emailCfg *config.EmailConfig) error {
	var count int64
	db.Model(&model.User{}).Where("username = ? OR email = ?", username, email).Count(&count)
	if count > 0 {
		return errors.New("用户名或邮箱已存在")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	token := generateToken(16)
	expire := time.Now().Add(24 * time.Hour)
	user := model.User{
		Username:          username,
		Email:             email,
		Password:          string(hash),
		Role:              "user",
		Status:            "pending",
		EmailVerified:     false,
		EmailVerifyToken:  token,
		EmailVerifyExpire: &expire,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	// 发送验证邮件（同步，便于调试）
	verifyURL := "http://localhost:5173/verify-email?token=" + token
	body := "请点击链接激活账号：<a href='" + verifyURL + "'>激活</a>，24小时内有效。"
	log.Printf("SendEmail config: %+v", emailCfg)
	err = SendEmail(emailCfg, user.Email, "F.Blog 邮箱验证", body)
	if err != nil {
		log.Printf("SendEmail error: %v", err)
		return errors.New("激活邮件发送失败: " + err.Error())
	}
	return nil
}

func LoginUser(db *gorm.DB, usernameOrEmail, password string) (string, *model.User, error) {
	var user model.User
	if err := db.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user).Error; err != nil {
		return "", nil, errors.New("用户不存在")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("密码错误")
	}
	if user.Status != "active" || !user.EmailVerified {
		return "", nil, errors.New("账号未激活，请先完成邮箱验证")
	}
	// 生成JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", nil, err
	}
	return tokenString, &user, nil
}

func ActivateEmail(c *gin.Context) {
	token := c.Query("token")
	var user model.User
	if err := database.DB.Where("email_verify_token = ? AND email_verify_expire > ?", token, time.Now()).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "激活链接无效或已过期"})
		return
	}
	user.Status = "active"
	user.EmailVerified = true
	user.EmailVerifyToken = ""
	user.EmailVerifyExpire = nil
	database.DB.Save(&user)
	c.JSON(200, gin.H{"message": "激活成功，请登录"})
}

// 密码加密
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// 密码校验
func CheckPassword(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
