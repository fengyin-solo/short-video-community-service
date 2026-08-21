package service

import (
	"sort"
	"time"

	"shortvideo/internal/model"
	"shortvideo/internal/store"
	"shortvideo/pkg/idgen"
)

// CreateVideo 创建视频（草稿）。
func (s *Service) CreateVideo(input model.Video) (*model.Video, error) {
	input.ID = idgen.Hex()
	input.Status = model.VideoStatusDraft
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if _, err := s.store.GetUser(input.AuthorID); err != nil {
		return nil, model.NewValidationError("author_id", "作者不存在")
	}
	if input.CategoryID != "" {
		if _, err := s.store.GetCategory(input.CategoryID); err != nil {
			return nil, model.NewValidationError("category_id", "分类不存在")
		}
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateVideo(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

// GetVideo 获取视频。
func (s *Service) GetVideo(id string) (*model.Video, error) {
	return s.store.GetVideo(id)
}

// ListVideos 分页列出视频。
func (s *Service) ListVideos(filter model.VideoFilter, page, size int) ([]*model.Video, int, error) {
	all := s.store.ListVideos()
	matched := make([]*model.Video, 0, len(all))
	for _, v := range all {
		if filter.Match(v) {
			matched = append(matched, v)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Video{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// SubmitVideoForReview 提交审核。
func (s *Service) SubmitVideoForReview(id string) (*model.Video, error) {
	return s.transitionVideo(id, model.VideoStatusReviewing)
}

// ApproveVideo 审核通过（发布）。
func (s *Service) ApproveVideo(id string) (*model.Video, error) {
	return s.transitionVideo(id, model.VideoStatusPublished)
}

// BanVideo 下架视频。
func (s *Service) BanVideo(id string) (*model.Video, error) {
	return s.transitionVideo(id, model.VideoStatusBanned)
}

func (s *Service) transitionVideo(id, to string) (*model.Video, error) {
	existing, err := s.store.GetVideo(id)
	if err != nil {
		return nil, err
	}
	if !model.CanVideoTransition(existing.Status, to) {
		return nil, store.ErrConflict
	}
	existing.Status = to
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateVideo(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// IncrementView 增加视频播放量。
func (s *Service) IncrementView(id string) (*model.Video, error) {
	existing, err := s.store.GetVideo(id)
	if err != nil {
		return nil, err
	}
	existing.ViewCount++
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateVideo(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteVideo 删除视频。
func (s *Service) DeleteVideo(id string) error {
	return s.store.DeleteVideo(id)
}
