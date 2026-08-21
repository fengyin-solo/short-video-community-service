package service

import (
	"sort"
	"time"

	"shortvideo/internal/model"
	"shortvideo/pkg/idgen"
)

// CreateUser 创建用户。
func (s *Service) CreateUser(input model.User) (*model.User, error) {
	input.ID = idgen.Hex()
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateUser(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

// GetUser 获取用户。
func (s *Service) GetUser(id string) (*model.User, error) {
	return s.store.GetUser(id)
}

// ListUsers 分页列出用户。
func (s *Service) ListUsers(filter model.UserFilter, page, size int) ([]*model.User, int, error) {
	all := s.store.ListUsers()
	matched := make([]*model.User, 0, len(all))
	for _, u := range all {
		if filter.Match(u) {
			matched = append(matched, u)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.User{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateUser 更新用户资料。
func (s *Service) UpdateUser(id string, input model.User) (*model.User, error) {
	existing, err := s.store.GetUser(id)
	if err != nil {
		return nil, err
	}
	existing.Username = input.Username
	existing.Nickname = input.Nickname
	existing.Avatar = input.Avatar
	existing.Bio = input.Bio
	existing.Status = input.Status
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateUser(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// BanUser 封禁用户。
func (s *Service) BanUser(id string) (*model.User, error) {
	existing, err := s.store.GetUser(id)
	if err != nil {
		return nil, err
	}
	existing.Status = model.UserBanned
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateUser(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// UnbanUser 解封用户。
func (s *Service) UnbanUser(id string) (*model.User, error) {
	existing, err := s.store.GetUser(id)
	if err != nil {
		return nil, err
	}
	existing.Status = model.UserActive
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateUser(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteUser 删除用户。
func (s *Service) DeleteUser(id string) error {
	return s.store.DeleteUser(id)
}
