package store

import "shortvideo/internal/model"

// CreateNotification 创建通知。
func (s *MemoryStore) CreateNotification(n *model.Notification) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notifications[n.ID] = n
	return nil
}

// GetNotification 按 ID 获取通知。
func (s *MemoryStore) GetNotification(id string) (*model.Notification, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.notifications[id]
	if !ok {
		return nil, ErrNotFound
	}
	return n, nil
}

// ListNotifications 列出全部通知。
func (s *MemoryStore) ListNotifications() []*model.Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Notification, 0, len(s.notifications))
	for _, n := range s.notifications {
		list = append(list, n)
	}
	return list
}

// UpdateNotification 更新通知。
func (s *MemoryStore) UpdateNotification(n *model.Notification) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.notifications[n.ID]; !ok {
		return ErrNotFound
	}
	s.notifications[n.ID] = n
	return nil
}

// DeleteNotification 删除通知。
func (s *MemoryStore) DeleteNotification(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.notifications[id]; !ok {
		return ErrNotFound
	}
	delete(s.notifications, id)
	return nil
}
