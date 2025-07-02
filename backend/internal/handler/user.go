package handler

import (
	"f_blog/backend/internal/database"
	"f_blog/backend/internal/service"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"f_blog/backend/internal/config"
	"f_blog/backend/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

var captchaStore = base64Captcha.DefaultMemStore

// 获取验证码图片
func GetCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverString(
		44,                               // height
		120,                              // width
		0,                                // noiseCount
		base64Captcha.OptionShowSineLine, // 干扰线
		4,                                // length
		"abcdefghjkmnpqrstuvwxyz23456789ABCDEFGHJKMNPQRSTUVWXYZ23456789", // 默认字符集
		nil, // 背景色默认
		nil, // fontsStorage 默认
		nil, // 字体颜色默认
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
	// 动态加载邮箱配置
	var emailCfgQQ = config.LoadEmailConfigQQ()
	var emailCfg163 = config.LoadEmailConfig163()
	// 根据邮箱后缀选择发件邮箱
	var emailCfg *config.EmailConfig
	if strings.HasSuffix(strings.ToLower(req.Email), "@qq.com") {
		emailCfg = emailCfgQQ
	} else if strings.HasSuffix(strings.ToLower(req.Email), "@163.com") {
		emailCfg = emailCfg163
	} else {
		emailCfg = emailCfgQQ // 默认用QQ
	}
	err := service.RegisterUser(database.DB, req.Username, req.Email, req.Password, emailCfg)
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
	token, user, err := service.LoginUser(database.DB, req.UsernameOrEmail, req.Password)
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
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
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
	updates := map[string]interface{}{
		"nickname": req.Nickname,
		"bio":      req.Bio,
		"gender":   req.Gender,
		"website":  req.Website,
	}
	if req.Birthday != "" {
		var t time.Time
		var err error
		if len(req.Birthday) == 10 {
			t, err = time.Parse("2006-01-02", req.Birthday)
		} else {
			t, err = time.Parse(time.RFC3339, req.Birthday)
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "生日格式错误"})
			return
		}
		updates["birthday"] = t
	}
	if err := database.DB.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
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
	if err := database.DB.Model(&model.User{}).Where("id = ?", userID).Update("avatar", avatarURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "头像更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"avatar": avatarURL})
}

func DeleteAvatar(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	if user.Avatar != "" {
		_ = os.Remove("." + user.Avatar)
	}
	if err := database.DB.Model(&model.User{}).Where("id = ?", userID).Update("avatar", "").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "头像删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "头像已删除"})
}

// 获取自己发布的文章
func ListMyArticles(c *gin.Context) {
	userID := c.GetUint("user_id")
	var articles []model.Article
	if err := database.DB.Where("author_id = ?", userID).Order("created_at desc").Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, articles)
}

// 获取自己发布的评论
func ListMyComments(c *gin.Context) {
	userID := c.GetUint("user_id")
	var comments []model.Comment
	if err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

// 获取自己点赞的内容
func ListMyLikes(c *gin.Context) {
	userID := c.GetUint("user_id")
	targetType := c.Query("type") // "article" or "comment"
	var likes []model.Like
	if err := database.DB.Where("user_id = ? AND target_type = ?", userID, targetType).Order("created_at desc").Find(&likes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, likes)
}

// 获取关注的人
func ListMyFollowings(c *gin.Context) {
	userID := c.GetUint("user_id")
	var follows []model.Follow
	if err := database.DB.Where("follower_id = ?", userID).Find(&follows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	// 批量查出被关注用户信息
	ids := make([]uint, 0, len(follows))
	for _, f := range follows {
		ids = append(ids, f.FollowingID)
	}
	var users []model.User
	database.DB.Where("id IN ?", ids).Find(&users)
	userMap := make(map[uint]model.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	type FollowWithUser struct {
		model.Follow
		FollowingUsername string `json:"following_username"`
		FollowingNickname string `json:"following_nickname"`
	}
	result := make([]FollowWithUser, 0, len(follows))
	for _, f := range follows {
		u := userMap[f.FollowingID]
		result = append(result, FollowWithUser{
			Follow:            f,
			FollowingUsername: u.Username,
			FollowingNickname: u.Nickname,
		})
	}
	c.JSON(http.StatusOK, result)
}

// 获取粉丝
func ListMyFollowers(c *gin.Context) {
	userID := c.GetUint("user_id")
	var follows []model.Follow
	if err := database.DB.Where("following_id = ?", userID).Find(&follows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	// 批量查出粉丝用户信息
	ids := make([]uint, 0, len(follows))
	for _, f := range follows {
		ids = append(ids, f.FollowerID)
	}
	var users []model.User
	database.DB.Where("id IN ?", ids).Find(&users)
	userMap := make(map[uint]model.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	type FollowerWithUser struct {
		model.Follow
		FollowerUsername string `json:"follower_username"`
		FollowerNickname string `json:"follower_nickname"`
	}
	result := make([]FollowerWithUser, 0, len(follows))
	for _, f := range follows {
		u := userMap[f.FollowerID]
		result = append(result, FollowerWithUser{
			Follow:           f,
			FollowerUsername: u.Username,
			FollowerNickname: u.Nickname,
		})
	}
	c.JSON(http.StatusOK, result)
}

// 关注
func Follow(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		FollowingID uint `json:"following_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if userID == req.FollowingID {
		c.JSON(400, gin.H{"error": "不能关注或取关自己"})
		return
	}
	follow := model.Follow{
		FollowerID:  userID,
		FollowingID: req.FollowingID,
	}
	// 防止重复关注
	if err := database.DB.Where(&follow).FirstOrCreate(&follow).Error; err != nil {
		c.JSON(500, gin.H{"error": "关注失败"})
		return
	}

	// 关注成功后，给被关注者发送通知（不通知自己）
	if userID != req.FollowingID {
		database.DB.Create(&model.Notification{
			UserID:  req.FollowingID,
			Type:    "follow",
			Title:   "你有新粉丝",
			Content: "有人关注了你",
			Data:    fmt.Sprintf(`{"follower_id":%d}`, userID),
		})
	}

	c.JSON(200, gin.H{"message": "关注成功"})
}

// 取消关注
func Unfollow(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		FollowingID uint `json:"following_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if userID == req.FollowingID {
		c.JSON(400, gin.H{"error": "不能关注或取关自己"})
		return
	}
	if err := database.DB.Where("follower_id = ? AND following_id = ?", userID, req.FollowingID).Delete(&model.Follow{}).Error; err != nil {
		c.JSON(500, gin.H{"error": "取消关注失败"})
		return
	}
	c.JSON(200, gin.H{"message": "取消关注成功"})
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

// 修改密码
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
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{"error": "用户不存在"})
		return
	}
	// 校验原密码
	if err := service.CheckPassword(user.Password, req.OldPassword); err != nil {
		c.JSON(400, gin.H{"error": "原密码错误"})
		return
	}
	// 加密新密码
	hash, err := service.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(500, gin.H{"error": "密码加密失败"})
		return
	}
	if err := database.DB.Model(&user).Update("password", hash).Error; err != nil {
		c.JSON(500, gin.H{"error": "修改失败"})
		return
	}
	c.JSON(200, gin.H{"message": "修改成功"})
}
