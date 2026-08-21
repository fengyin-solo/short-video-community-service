package store

import "shortvideo/internal/model"

// CreateVideo 创建视频。
func (s *MemoryStore) CreateVideo(v *model.Video) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.videos[v.ID] = v
	return nil
}

// GetVideo 按 ID 获取视频。
func (s *MemoryStore) GetVideo(id string) (*model.Video, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.videos[id]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

// ListVideos 列出全部视频。
func (s *MemoryStore) ListVideos() []*model.Video {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Video, 0, len(s.videos))
	for _, v := range s.videos {
		list = append(list, v)
	}
	return list
}

// UpdateVideo 更新视频。
func (s *MemoryStore) UpdateVideo(v *model.Video) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.videos[v.ID]; !ok {
		return ErrNotFound
	}
	s.videos[v.ID] = v
	return nil
}

// DeleteVideo 删除视频。
func (s *MemoryStore) DeleteVideo(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.videos[id]; !ok {
		return ErrNotFound
	}
	delete(s.videos, id)
	return nil
}
