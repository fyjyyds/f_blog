package service

import (
	"f_blog/backend/internal/cache"
	"f_blog/backend/internal/model"
	"f_blog/backend/internal/repository"
	"encoding/json"
)

type TagService struct {
	repo *repository.TagRepository
}

func NewTagService(repo *repository.TagRepository) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) List() ([]model.Tag, error) {
	// 尝试从缓存读取
	if cached, err := cache.Get(cache.KeyTagsList); err == nil {
		var tags []model.Tag
		if json.Unmarshal([]byte(cached), &tags) == nil {
			return tags, nil
		}
	}

	// 缓存未命中，查数据库
	tags, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if data, err := json.Marshal(tags); err == nil {
		cache.Set(cache.KeyTagsList, string(data), cache.TTLTags)
	}

	return tags, nil
}

func (s *TagService) AdminList() ([]model.Tag, error) {
	return s.repo.FindAll()
}

func (s *TagService) Create(name, color string) (*model.Tag, error) {
	if color == "" {
		color = "#667eea"
	}
	tag := &model.Tag{
		Name:  name,
		Color: color,
	}
	err := s.repo.Create(tag)
	if err == nil {
		cache.Delete(cache.KeyTagsList)
	}
	return tag, err
}

func (s *TagService) Update(id uint, name, color string) (*model.Tag, error) {
	tag, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	updates := map[string]interface{}{
		"name":  name,
		"color": color,
	}
	err = s.repo.Update(tag, updates)
	if err == nil {
		cache.Delete(cache.KeyTagsList)
	}
	return tag, err
}

func (s *TagService) Delete(id uint) error {
	err := s.repo.Delete(id)
	if err == nil {
		cache.Delete(cache.KeyTagsList)
	}
	return err
}
