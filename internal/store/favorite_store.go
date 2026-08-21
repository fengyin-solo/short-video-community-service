package store

import "shortvideo/internal/model"

// CreateFavorite 创建收藏记录。
func (s *MemoryStore) CreateFavorite(f *model.Favorite) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.favorites {
		if exist.UserID == f.UserID && exist.VideoID == f.VideoID {
			return ErrConflict
		}
	}
	s.favorites[f.ID] = f
	return nil
}

// GetFavorite 按 ID 获取收藏记录。
func (s *MemoryStore) GetFavorite(id string) (*model.Favorite, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.favorites[id]
	if !ok {
		return nil, ErrNotFound
	}
	return f, nil
}

// GetFavoriteByPair 按用户+视频获取收藏记录。
func (s *MemoryStore) GetFavoriteByPair(userID, videoID string) (*model.Favorite, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, f := range s.favorites {
		if f.UserID == userID && f.VideoID == videoID {
			return f, nil
		}
	}
	return nil, ErrNotFound
}

// ListFavorites 列出全部收藏记录。
func (s *MemoryStore) ListFavorites() []*model.Favorite {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Favorite, 0, len(s.favorites))
	for _, f := range s.favorites {
		list = append(list, f)
	}
	return list
}

// DeleteFavorite 删除收藏记录。
func (s *MemoryStore) DeleteFavorite(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.favorites[id]; !ok {
		return ErrNotFound
	}
	delete(s.favorites, id)
	return nil
}
