package service

import (
	"sort"
	"time"

	"shortvideo/internal/model"
	"shortvideo/pkg/idgen"
)

// CreateTag 创建标签。
func (s *Service) CreateTag(input model.Tag) (*model.Tag, error) {
	input.ID = idgen.Hex()
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateTag(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

// GetTag 获取标签。
func (s *Service) GetTag(id string) (*model.Tag, error) {
	return s.store.GetTag(id)
}

// ListTags 分页列出标签。
func (s *Service) ListTags(page, size int) ([]*model.Tag, int, error) {
	all := s.store.ListTags()
	sort.Slice(all, func(i, j int) bool { return all[i].VideoCount > all[j].VideoCount })
	total := len(all)
	start := (page - 1) * size
	if start >= total {
		return []*model.Tag{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

// DeleteTag 删除标签。
func (s *Service) DeleteTag(id string) error {
	return s.store.DeleteTag(id)
}
