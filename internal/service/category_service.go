package service

import (
	"sort"
	"time"

	"shortvideo/internal/model"
	"shortvideo/pkg/idgen"
)

// CreateCategory 创建分类。
func (s *Service) CreateCategory(input model.Category) (*model.Category, error) {
	input.ID = idgen.Hex()
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateCategory(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

// GetCategory 获取分类。
func (s *Service) GetCategory(id string) (*model.Category, error) {
	return s.store.GetCategory(id)
}

// ListCategories 分页列出分类。
func (s *Service) ListCategories(page, size int) ([]*model.Category, int, error) {
	all := s.store.ListCategories()
	sort.Slice(all, func(i, j int) bool { return all[i].Sort < all[j].Sort })
	total := len(all)
	start := (page - 1) * size
	if start >= total {
		return []*model.Category{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

// UpdateCategory 更新分类。
func (s *Service) UpdateCategory(id string, input model.Category) (*model.Category, error) {
	existing, err := s.store.GetCategory(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Sort = input.Sort
	existing.Status = input.Status
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateCategory(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteCategory 删除分类。
func (s *Service) DeleteCategory(id string) error {
	return s.store.DeleteCategory(id)
}
