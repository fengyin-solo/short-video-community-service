package store

import "shortvideo/internal/model"

// CreateReport 创建举报记录。
func (s *MemoryStore) CreateReport(r *model.Report) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reports[r.ID] = r
	return nil
}

// GetReport 按 ID 获取举报记录。
func (s *MemoryStore) GetReport(id string) (*model.Report, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.reports[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

// ListReports 列出全部举报记录。
func (s *MemoryStore) ListReports() []*model.Report {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Report, 0, len(s.reports))
	for _, r := range s.reports {
		list = append(list, r)
	}
	return list
}

// UpdateReport 更新举报记录。
func (s *MemoryStore) UpdateReport(r *model.Report) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.reports[r.ID]; !ok {
		return ErrNotFound
	}
	s.reports[r.ID] = r
	return nil
}

// DeleteReport 删除举报记录。
func (s *MemoryStore) DeleteReport(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.reports[id]; !ok {
		return ErrNotFound
	}
	delete(s.reports, id)
	return nil
}
