package service

import (
	"f_blog/backend/internal/repository"
	"strings"
)

type SettingService struct {
	repo *repository.SettingRepository
}

func NewSettingService(repo *repository.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

func (s *SettingService) GetGrouped() (map[string]interface{}, error) {
	settings, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	basic := make(map[string]interface{})
	email := make(map[string]interface{})
	security := make(map[string]interface{})
	content := make(map[string]interface{})

	for _, setting := range settings {
		if strings.HasPrefix(setting.Key, "email_") {
			key := strings.TrimPrefix(setting.Key, "email_")
			email[key] = setting.Value
		} else if strings.HasPrefix(setting.Key, "security_") {
			key := strings.TrimPrefix(setting.Key, "security_")
			security[key] = setting.Value
		} else if strings.HasPrefix(setting.Key, "content_") {
			key := strings.TrimPrefix(setting.Key, "content_")
			content[key] = setting.Value
		} else {
			basic[setting.Key] = setting.Value
		}
	}

	// 默认值
	if len(basic) == 0 {
		basic = map[string]interface{}{
			"siteName":        "F_Blog",
			"siteDescription": "一个现代化的博客系统",
			"siteKeywords":    "",
			"siteLogo":        "",
			"icp":             "",
		}
	}
	if len(email) == 0 {
		email = map[string]interface{}{
			"smtpHost":      "smtp.example.com",
			"smtpPort":      "587",
			"emailUser":     "admin@example.com",
			"emailPassword": "",
			"senderName":    "F_Blog",
		}
	}
	if len(security) == 0 {
		security = map[string]interface{}{
			"minPasswordLength": "8",
			"maxLoginAttempts":  "5",
			"lockoutDuration":   "30",
			"jwtExpireHours":    "24",
			"enableCaptcha":     "false",
			"enableTwoFactor":   "false",
		}
	}
	if len(content) == 0 {
		content = map[string]interface{}{
			"articlesPerPage":        "10",
			"summaryLength":          "200",
			"commentModeration":      "manual",
			"allowAnonymousComments": "false",
			"enableCommentLikes":     "true",
			"enableArticleLikes":     "true",
		}
	}

	return map[string]interface{}{
		"basic":    basic,
		"email":    email,
		"security": security,
		"content":  content,
	}, nil
}

func (s *SettingService) Update(basic, email, security, content map[string]interface{}) error {
	if basic != nil {
		for key, value := range basic {
			if err := s.repo.Upsert(key, value, "string"); err != nil {
				return err
			}
		}
	}
	if email != nil {
		for key, value := range email {
			if err := s.repo.Upsert("email_"+key, value, "string"); err != nil {
				return err
			}
		}
	}
	if security != nil {
		for key, value := range security {
			if err := s.repo.Upsert("security_"+key, value, "string"); err != nil {
				return err
			}
		}
	}
	if content != nil {
		for key, value := range content {
			if err := s.repo.Upsert("content_"+key, value, "string"); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *SettingService) Reset() error {
	return s.repo.ResetAll()
}
