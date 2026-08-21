package store

import "shortvideo/internal/model"

// CreateLike 创建点赞记录。
func (s *MemoryStore) CreateLike(l *model.Like) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.likes {
		if exist.UserID == l.UserID && exist.VideoID == l.VideoID {
			return ErrConflict
		}
	}
	s.likes[l.ID] = l
	return nil
}

// GetLike 按 ID 获取点赞记录。
func (s *MemoryStore) GetLike(id string) (*model.Like, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.likes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return l, nil
}

// GetLikeByPair 按用户+视频获取点赞记录。
func (s *MemoryStore) GetLikeByPair(userID, videoID string) (*model.Like, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, l := range s.likes {
		if l.UserID == userID && l.VideoID == videoID {
			return l, nil
		}
	}
	return nil, ErrNotFound
}

// ListLikes 列出全部点赞记录。
func (s *MemoryStore) ListLikes() []*model.Like {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Like, 0, len(s.likes))
	for _, l := range s.likes {
		list = append(list, l)
	}
	return list
}

// DeleteLike 删除点赞记录。
func (s *MemoryStore) DeleteLike(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.likes[id]; !ok {
		return ErrNotFound
	}
	delete(s.likes, id)
	return nil
}
