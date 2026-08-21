package store

import "sync"

type FeedRecipientStore struct {
	mu      sync.RWMutex
	userIDs []string
}

func NewFeedRecipientStore(userIDs []string) *FeedRecipientStore {
	return &FeedRecipientStore{userIDs: append([]string(nil), userIDs...)}
}

func (s *FeedRecipientStore) Snapshot() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.userIDs
}

func (s *FeedRecipientStore) ReplaceFirst(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.userIDs[0] = userID
}
