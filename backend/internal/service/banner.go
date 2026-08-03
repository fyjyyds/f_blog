package service

import (
	"encoding/json"
	"f_blog/backend/internal/cache"
	"f_blog/backend/internal/model"
	"f_blog/backend/internal/repository"
	"f_blog/backend/internal/util"
)

type BannerService struct {
	repo *repository.BannerRepository
}

func NewBannerService(repo *repository.BannerRepository) *BannerService {
	return &BannerService{repo: repo}
}

func (s *BannerService) ListActive() ([]model.Banner, error) {
	// 尝试从缓存读取
	if cached, err := cache.Get(cache.KeyBannersActive); err == nil {
		var banners []model.Banner
		if json.Unmarshal([]byte(cached), &banners) == nil {
			return s.processURLs(banners), nil
		}
	}

	banners, err := s.repo.FindActive()
	if err != nil {
		return nil, err
	}

	// 写入缓存（存原始数据，不含拼接后的 URL）
	if data, err := json.Marshal(banners); err == nil {
		cache.Set(cache.KeyBannersActive, string(data), cache.TTLBanners)
	}

	return s.processURLs(banners), nil
}

func (s *BannerService) AdminList() ([]model.Banner, error) {
	banners, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return s.processURLs(banners), nil
}

func (s *BannerService) Create(title, image, link string, sortOrder int, status string) (*model.Banner, error) {
	if status == "" {
		status = "active"
	}
	// 存相对路径
	banner := &model.Banner{
		Title:     title,
		Image:     util.RelativePath(image),
		Link:      link,
		SortOrder: sortOrder,
		Status:    status,
	}
	err := s.repo.Create(banner)
	if err == nil {
		cache.Delete(cache.KeyBannersActive)
	}
	// 返回时拼接完整 URL
	banner.Image = util.StaticURL(banner.Image)
	return banner, err
}

func (s *BannerService) Update(id uint, title, image, link string, sortOrder int, status string) (*model.Banner, error) {
	banner, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	updates := map[string]interface{}{
		"title":      title,
		"image":      util.RelativePath(image),
		"link":       link,
		"sort_order": sortOrder,
		"status":     status,
	}
	err = s.repo.Update(banner, updates)
	if err == nil {
		cache.Delete(cache.KeyBannersActive)
		banner.Image = util.StaticURL(util.RelativePath(image))
	}
	return banner, err
}

func (s *BannerService) Delete(id uint) error {
	err := s.repo.Delete(id)
	if err == nil {
		cache.Delete(cache.KeyBannersActive)
	}
	return err
}

// processURLs 为 banner 列表拼接完整的图片 URL
func (s *BannerService) processURLs(banners []model.Banner) []model.Banner {
	for i := range banners {
		banners[i].Image = util.StaticURL(banners[i].Image)
	}
	return banners
}
