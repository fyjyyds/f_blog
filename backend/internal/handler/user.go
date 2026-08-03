package handler

import (
	"f_blog/backend/internal/config"
	"f_blog/backend/internal/service"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

var captchaStore = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverString(
		44, 120, 0,
		base64Captcha.OptionShowSineLine,
		4,
		"abcdefghjkmnpqrstuvwxyz23456789ABCDEFGHJKMNPQRSTUVWXYZ23456789",
		nil, nil, nil,
	)
	captcha := base64Captcha.NewCaptcha(driver, captchaStore)
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		c.JSON(500, gin.H{"error": "生成验证码失败"})
		return
	}
	c.JSON(200, gin.H{
		"captcha_id":    id,
		"captcha_image": b64s,
	})
}

func verifyCaptcha(captchaId, value string) bool {
	return captchaStore.Verify(captchaId, value, true)
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=20"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=20"`
	CaptchaID string `json:"captcha_id"`
	Captcha   string `json:"captcha"`
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	CaptchaID       string `json:"captcha_id"`
	Captcha         string `json:"captcha"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.CaptchaID == "" || req.Captcha == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码不能为空"})
		return
	}
	if !verifyCaptcha(req.CaptchaID, req.Captcha) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	var emailCfg *config.EmailConfig
	if strings.HasSuffix(strings.ToLower(req.Email), "@163.com") {
		emailCfg = config.LoadEmailConfig163()
	} else {
		emailCfg = config.LoadEmailConfigQQ()
	}

	err := service.App.User.Register(req.Username, req.Email, req.Password, emailCfg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "注册成功，请前往邮箱激活账号"})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.CaptchaID == "" || req.Captcha == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码不能为空"})
		return
	}
	if !verifyCaptcha(req.CaptchaID, req.Captcha) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	token, user, err := service.App.User.Login(req.UsernameOrEmail, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
		},
	})
}

func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := service.App.User.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Nickname string `json:"nickname"`
		Bio      string `json:"bio"`
		Gender   string `json:"gender"`
		Birthday string `json:"birthday"`
		Website  string `json:"website"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.App.User.UpdateProfile(userID, req.Nickname, req.Bio, req.Gender, req.Birthday, req.Website); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func UploadAvatar(c *gin.Context) {
	userID := c.GetUint("user_id")
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择头像文件"})
		return
	}
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持jpg/jpeg/png/gif格式"})
		return
	}
	if file.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过2MB"})
		return
	}
	saveDir := "static/avatar"
	_ = os.MkdirAll(saveDir, os.ModePerm)
	filename := "avatar_" + strconv.FormatUint(uint64(userID), 10) + ext
	savePath := filepath.Join(saveDir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}
	avatarURL := "/static/avatar/" + filename
	if err := service.App.User.UpdateAvatar(userID, avatarURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "头像更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"avatar": avatarURL})
}

func DeleteAvatar(c *gin.Context) {
	userID := c.GetUint("user_id")
	oldAvatar, err := service.App.User.DeleteAvatar(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if oldAvatar != "" {
		_ = os.Remove("." + oldAvatar)
	}
	c.JSON(http.StatusOK, gin.H{"message": "头像已删除"})
}

func ListMyArticles(c *gin.Context) {
	userID := c.GetUint("user_id")
	articles, err := service.App.Article.ListByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, articles)
}

func ListMyComments(c *gin.Context) {
	userID := c.GetUint("user_id")
	comments, err := service.App.Comment.ListByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func ListMyLikes(c *gin.Context) {
	userID := c.GetUint("user_id")
	targetType := c.Query("type")
	likes, err := service.App.Like.ListByUserID(userID, targetType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, likes)
}

func ListMyFollowings(c *gin.Context) {
	userID := c.GetUint("user_id")
	result, err := service.App.Follow.ListFollowings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func ListMyFollowers(c *gin.Context) {
	userID := c.GetUint("user_id")
	result, err := service.App.Follow.ListFollowers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func Follow(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		FollowingID uint `json:"following_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := service.App.Follow.Follow(userID, req.FollowingID); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "关注成功"})
}

func Unfollow(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		FollowingID uint `json:"following_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := service.App.Follow.Unfollow(userID, req.FollowingID); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "取消关注成功"})
}

func ActivateEmail(c *gin.Context) {
	token := c.Query("token")
	_, err := service.App.User.ActivateEmail(token)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "激活成功，请登录"})
}

func ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	if err := service.App.User.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "修改成功"})
}

// getIntFromString 辅助函数
func getIntFromString(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return defaultValue
}
