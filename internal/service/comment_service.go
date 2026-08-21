package service

import (
	"sort"
	"time"

	"shortvideo/internal/model"
	"shortvideo/pkg/idgen"
)

// CreateComment 发表评论（联动视频评论数 +1）。
func (s *Service) CreateComment(input model.Comment) (*model.Comment, error) {
	input.ID = idgen.Hex()
	input.Status = model.CommentNormal
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if _, err := s.store.GetVideo(input.VideoID); err != nil {
		return nil, model.NewValidationError("video_id", "视频不存在")
	}
	if _, err := s.store.GetUser(input.UserID); err != nil {
		return nil, model.NewValidationError("user_id", "用户不存在")
	}
	if input.ParentID != "" {
		if _, err := s.store.GetComment(input.ParentID); err != nil {
			return nil, model.NewValidationError("parent_id", "父评论不存在")
		}
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateComment(&input); err != nil {
		return nil, err
	}
	if v, err := s.store.GetVideo(input.VideoID); err == nil {
		v.CommentCount++
		v.UpdatedAt = now
		_ = s.store.UpdateVideo(v)
	}
	return &input, nil
}

// GetComment 获取评论。
func (s *Service) GetComment(id string) (*model.Comment, error) {
	return s.store.GetComment(id)
}

// ListComments 分页列出评论。
func (s *Service) ListComments(filter model.CommentFilter, page, size int) ([]*model.Comment, int, error) {
	all := s.store.ListComments()
	matched := make([]*model.Comment, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Comment{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// DeleteComment 删除评论（软删，联动视频评论数 -1）。
func (s *Service) DeleteComment(id string) (*model.Comment, error) {
	existing, err := s.store.GetComment(id)
	if err != nil {
		return nil, err
	}
	if existing.Status == model.CommentDeleted {
		return existing, nil
	}
	existing.Status = model.CommentDeleted
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateComment(existing); err != nil {
		return nil, err
	}
	if v, err := s.store.GetVideo(existing.VideoID); err == nil {
		if v.CommentCount > 0 {
			v.CommentCount--
			v.UpdatedAt = time.Now()
			_ = s.store.UpdateVideo(v)
		}
	}
	return existing, nil
}
