package store

import "shortvideo/internal/model"

// CreatePlaylist 创建播单。
func (s *MemoryStore) CreatePlaylist(p *model.Playlist) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.playlists[p.ID] = p
	return nil
}

// GetPlaylist 按 ID 获取播单。
func (s *MemoryStore) GetPlaylist(id string) (*model.Playlist, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.playlists[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

// ListPlaylists 列出全部播单。
func (s *MemoryStore) ListPlaylists() []*model.Playlist {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Playlist, 0, len(s.playlists))
	for _, p := range s.playlists {
		list = append(list, p)
	}
	return list
}

// UpdatePlaylist 更新播单。
func (s *MemoryStore) UpdatePlaylist(p *model.Playlist) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.playlists[p.ID]; !ok {
		return ErrNotFound
	}
	s.playlists[p.ID] = p
	return nil
}

// DeletePlaylist 删除播单。
func (s *MemoryStore) DeletePlaylist(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.playlists[id]; !ok {
		return ErrNotFound
	}
	delete(s.playlists, id)
	return nil
}
