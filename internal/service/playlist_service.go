package service

import (
	"sort"
	"time"

	"shortvideo/internal/model"
	"shortvideo/pkg/idgen"
)

// CreatePlaylist 创建播单。
func (s *Service) CreatePlaylist(input model.Playlist) (*model.Playlist, error) {
	input.ID = idgen.Hex()
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if _, err := s.store.GetUser(input.UserID); err != nil {
		return nil, model.NewValidationError("user_id", "用户不存在")
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreatePlaylist(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

// GetPlaylist 获取播单。
func (s *Service) GetPlaylist(id string) (*model.Playlist, error) {
	return s.store.GetPlaylist(id)
}

// ListPlaylists 分页列出播单。
func (s *Service) ListPlaylists(filter model.PlaylistFilter, page, size int) ([]*model.Playlist, int, error) {
	all := s.store.ListPlaylists()
	matched := make([]*model.Playlist, 0, len(all))
	for _, p := range all {
		if filter.Match(p) {
			matched = append(matched, p)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Playlist{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdatePlaylist 更新播单基本信息。
func (s *Service) UpdatePlaylist(id string, input model.Playlist) (*model.Playlist, error) {
	existing, err := s.store.GetPlaylist(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Description = input.Description
	existing.Status = input.Status
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdatePlaylist(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// AddVideoToPlaylist 向播单添加视频。
func (s *Service) AddVideoToPlaylist(id, videoID string) (*model.Playlist, error) {
	existing, err := s.store.GetPlaylist(id)
	if err != nil {
		return nil, err
	}
	if _, err := s.store.GetVideo(videoID); err != nil {
		return nil, model.NewValidationError("video_id", "视频不存在")
	}
	for _, vid := range existing.VideoIDs {
		if vid == videoID {
			return existing, nil
		}
	}
	existing.VideoIDs = append(existing.VideoIDs, videoID)
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdatePlaylist(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// RemoveVideoFromPlaylist 从播单移除视频。
func (s *Service) RemoveVideoFromPlaylist(id, videoID string) (*model.Playlist, error) {
	existing, err := s.store.GetPlaylist(id)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(existing.VideoIDs))
	for _, vid := range existing.VideoIDs {
		if vid != videoID {
			ids = append(ids, vid)
		}
	}
	existing.VideoIDs = ids
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdatePlaylist(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeletePlaylist 删除播单。
func (s *Service) DeletePlaylist(id string) error {
	return s.store.DeletePlaylist(id)
}
