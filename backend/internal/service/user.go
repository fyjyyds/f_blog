package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"f_blog/backend/internal/cache"
	"f_blog/backend/internal/config"
	"f_blog/backend/internal/model"
	"f_blog/backend/internal/repository"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func generateToken(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *UserService) Register(username, email, password string, emailCfg *config.EmailConfig) error {
	exists, err := s.repo.ExistsByUsernameOrEmail(username, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("用户名或邮箱已存在")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	token := generateToken(16)
	expire := time.Now().Add(24 * time.Hour)
	user := &model.User{
		Username:          username,
		Email:             email,
		Password:          string(hash),
		Role:              "user",
		Status:            "pending",
		EmailVerified:     false,
		EmailVerifyToken:  token,
		EmailVerifyExpire: &expire,
	}
	if err := s.repo.Create(user); err != nil {
		return err
	}

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

func (s *UserService) Login(usernameOrEmail, password string) (string, *model.User, error) {
	user, err := s.repo.FindByUsernameOrEmail(usernameOrEmail)
	if err != nil {
		return "", nil, errors.New("用户不存在")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("密码错误")
	}
	if user.Status != "active" || !user.EmailVerified {
		return "", nil, errors.New("账号未激活，请先完成邮箱验证")
	}

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
	return tokenString, user, nil
}

func (s *UserService) GetProfile(userID uint) (*model.User, error) {
	// 尝试从缓存读取
	cacheKey := cache.KeyUserProfile(userID)
	if cached, err := cache.Get(cacheKey); err == nil {
		var user model.User
		if json.Unmarshal([]byte(cached), &user) == nil {
			return &user, nil
		}
	}

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if data, err := json.Marshal(user); err == nil {
		cache.Set(cacheKey, string(data), cache.TTLUserProfile)
	}

	return user, nil
}

func (s *UserService) UpdateProfile(userID uint, nickname, bio, gender, birthday, website string) error {
	updates := map[string]interface{}{
		"nickname": nickname,
		"bio":      bio,
		"gender":   gender,
		"website":  website,
	}
	if birthday != "" {
		var t time.Time
		var err error
		if len(birthday) == 10 {
			t, err = time.Parse("2006-01-02", birthday)
		} else {
			t, err = time.Parse(time.RFC3339, birthday)
		}
		if err != nil {
			return errors.New("生日格式错误")
		}
		updates["birthday"] = t
	}
	err := s.repo.UpdateProfile(userID, updates)
	if err == nil {
		cache.Delete(cache.KeyUserProfile(userID))
	}
	return err
}

func (s *UserService) UpdateAvatar(userID uint, avatarURL string) error {
	err := s.repo.UpdateProfile(userID, map[string]interface{}{"avatar": avatarURL})
	if err == nil {
		cache.Delete(cache.KeyUserProfile(userID))
	}
	return err
}

func (s *UserService) DeleteAvatar(userID uint) (string, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return "", errors.New("用户不存在")
	}
	oldAvatar := user.Avatar
	if err := s.repo.UpdateProfile(userID, map[string]interface{}{"avatar": ""}); err != nil {
		return "", err
	}
	cache.Delete(cache.KeyUserProfile(userID))
	return oldAvatar, nil
}

func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("原密码错误")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}
	return s.repo.UpdatePassword(userID, string(hash))
}

func (s *UserService) ActivateEmail(token string) (*model.User, error) {
	user, err := s.repo.FindByVerifyToken(token, time.Now())
	if err != nil {
		return nil, errors.New("激活链接无效或已过期")
	}
	if err := s.repo.ActivateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Count 统计用户总数
func (s *UserService) Count() (int64, error) {
	return s.repo.Count()
}

// Admin operations

func (s *UserService) AdminList(page, pageSize int, search, status string) ([]model.User, int64, error) {
	return s.repo.List(page, pageSize, search, status)
}

func (s *UserService) AdminUpdate(userID uint, username, email, nickname, role, status string) (*model.User, error) {
	_, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	updates := map[string]interface{}{
		"username": username,
		"email":    email,
		"nickname": nickname,
		"role":     role,
		"status":   status,
	}
	if err := s.repo.Update(userID, updates); err != nil {
		return nil, err
	}
	cache.Delete(cache.KeyUserProfile(userID))
	return s.repo.FindByID(userID)
}

func (s *UserService) AdminBan(userID uint) error {
	_, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}
	err = s.repo.UpdateStatus(userID, "banned")
	if err == nil {
		cache.Delete(cache.KeyUserProfile(userID))
	}
	return err
}

func (s *UserService) AdminUnban(userID uint) error {
	_, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}
	err = s.repo.UpdateStatus(userID, "active")
	if err == nil {
		cache.Delete(cache.KeyUserProfile(userID))
	}
	return err
}

func (s *UserService) AdminDelete(userID uint) error {
	err := s.repo.Delete(userID)
	if err == nil {
		cache.Delete(cache.KeyUserProfile(userID))
	}
	return err
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPassword(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
