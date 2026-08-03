package service

import (
	"f_blog/backend/internal/cache"
	"f_blog/backend/internal/model"
	"f_blog/backend/internal/repository"
	"encoding/json"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) List() ([]model.Category, error) {
	// 尝试从缓存读取
	if cached, err := cache.Get(cache.KeyCategoriesList); err == nil {
		var categories []model.Category
		if json.Unmarshal([]byte(cached), &categories) == nil {
			return categories, nil
		}
	}

	// 缓存未命中，查数据库
	categories, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if data, err := json.Marshal(categories); err == nil {
		cache.Set(cache.KeyCategoriesList, string(data), cache.TTLCategories)
	}

	return categories, nil
}

func (s *CategoryService) AdminList() ([]model.Category, error) {
	// Admin 列表不缓存，直接查库
	return s.repo.FindAll()
}

func (s *CategoryService) Create(name, description string, sortOrder int, status string) (*model.Category, error) {
	if status == "" {
		status = "active"
	}
	category := &model.Category{
		Name:        name,
		Description: description,
		SortOrder:   sortOrder,
		Status:      status,
	}
	err := s.repo.Create(category)
	if err == nil {
		cache.Delete(cache.KeyCategoriesList) // 失效缓存
	}
	return category, err
}

func (s *CategoryService) Update(id uint, name, description string, sortOrder int, status string) (*model.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	updates := map[string]interface{}{
		"name":        name,
		"description": description,
		"sort_order":  sortOrder,
		"status":      status,
	}
	err = s.repo.Update(category, updates)
	if err == nil {
		cache.Delete(cache.KeyCategoriesList)
	}
	return category, err
}

func (s *CategoryService) Delete(id uint) error {
	err := s.repo.Delete(id)
	if err == nil {
		cache.Delete(cache.KeyCategoriesList)
	}
	return err
}
