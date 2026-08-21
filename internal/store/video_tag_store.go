package store

import "shortvideo/internal/model"

// CreateVideoTag 创建视频标签关联。
func (s *MemoryStore) CreateVideoTag(v *model.VideoTag) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.videoTags {
		if exist.VideoID == v.VideoID && exist.TagID == v.TagID {
			return ErrConflict
		}
	}
	s.videoTags[v.ID] = v
	return nil
}

// GetVideoTag 按 ID 获取视频标签关联。
func (s *MemoryStore) GetVideoTag(id string) (*model.VideoTag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.videoTags[id]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

// GetVideoTagByPair 按视频+标签获取关联。
func (s *MemoryStore) GetVideoTagByPair(videoID, tagID string) (*model.VideoTag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.videoTags {
		if v.VideoID == videoID && v.TagID == tagID {
			return v, nil
		}
	}
	return nil, ErrNotFound
}

// ListVideoTags 列出全部视频标签关联。
func (s *MemoryStore) ListVideoTags() []*model.VideoTag {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.VideoTag, 0, len(s.videoTags))
	for _, v := range s.videoTags {
		list = append(list, v)
	}
	return list
}

// DeleteVideoTag 删除视频标签关联。
func (s *MemoryStore) DeleteVideoTag(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.videoTags[id]; !ok {
		return ErrNotFound
	}
	delete(s.videoTags, id)
	return nil
}
