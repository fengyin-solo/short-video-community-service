package store

import "shortvideo/internal/model"

// CreateFollow 创建关注关系。
func (s *MemoryStore) CreateFollow(f *model.Follow) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.follows {
		if exist.FollowerID == f.FollowerID && exist.FolloweeID == f.FolloweeID {
			return ErrConflict
		}
	}
	s.follows[f.ID] = f
	return nil
}

// GetFollow 按 ID 获取关注关系。
func (s *MemoryStore) GetFollow(id string) (*model.Follow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.follows[id]
	if !ok {
		return nil, ErrNotFound
	}
	return f, nil
}

// GetFollowByPair 按关注者+被关注者获取关注关系。
func (s *MemoryStore) GetFollowByPair(followerID, followeeID string) (*model.Follow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, f := range s.follows {
		if f.FollowerID == followerID && f.FolloweeID == followeeID {
			return f, nil
		}
	}
	return nil, ErrNotFound
}

// ListFollows 列出全部关注关系。
func (s *MemoryStore) ListFollows() []*model.Follow {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Follow, 0, len(s.follows))
	for _, f := range s.follows {
		list = append(list, f)
	}
	return list
}

// DeleteFollow 删除关注关系。
func (s *MemoryStore) DeleteFollow(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.follows[id]; !ok {
		return ErrNotFound
	}
	delete(s.follows, id)
	return nil
}
