package repository

import (
	"f_blog/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) FindByUsernameOrEmail(usernameOrEmail string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user).Error
	return &user, err
}

func (r *UserRepository) ExistsByUsernameOrEmail(username, email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("username = ? OR email = ?", username, email).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) Create(user *model.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return r.db.Create(user).Error
}

func (r *UserRepository) UpdateProfile(userID uint, updates map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (r *UserRepository) UpdateAvatar(userID uint, avatarURL string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("avatar", avatarURL).Error
}

func (r *UserRepository) UpdatePassword(userID uint, hashedPassword string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}

func (r *UserRepository) ActivateUser(user *model.User) error {
	user.Status = "active"
	user.EmailVerified = true
	user.EmailVerifyToken = ""
	user.EmailVerifyExpire = nil
	return r.db.Save(user).Error
}

func (r *UserRepository) FindByVerifyToken(token string, now time.Time) (*model.User, error) {
	var user model.User
	err := r.db.Where("email_verify_token = ? AND email_verify_expire > ?", token, now).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByIDs(ids []uint) ([]model.User, error) {
	var users []model.User
	err := r.db.Where("id IN ?", ids).Find(&users).Error
	return users, err
}

// Count 统计用户总数
func (r *UserRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	return count, err
}

// Admin operations

func (r *UserRepository) List(page, pageSize int, search, status string) ([]model.User, int64, error) {
	db := r.db.Model(&model.User{})

	if search != "" {
		db = db.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}

	var total int64
	db.Count(&total)

	var users []model.User
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at desc").Find(&users).Error
	return users, total, err
}

func (r *UserRepository) Update(userID uint, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return r.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (r *UserRepository) UpdateStatus(userID uint, status string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("status", status).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}
