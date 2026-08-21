package store

import "shortvideo/internal/model"

// CreateComment 创建评论。
func (s *MemoryStore) CreateComment(c *model.Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.comments[c.ID] = c
	return nil
}

// GetComment 按 ID 获取评论。
func (s *MemoryStore) GetComment(id string) (*model.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.comments[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

// ListComments 列出全部评论。
func (s *MemoryStore) ListComments() []*model.Comment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Comment, 0, len(s.comments))
	for _, c := range s.comments {
		list = append(list, c)
	}
	return list
}

// UpdateComment 更新评论。
func (s *MemoryStore) UpdateComment(c *model.Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.comments[c.ID]; !ok {
		return ErrNotFound
	}
	s.comments[c.ID] = c
	return nil
}

// DeleteComment 删除评论。
func (s *MemoryStore) DeleteComment(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.comments[id]; !ok {
		return ErrNotFound
	}
	delete(s.comments, id)
	return nil
}
