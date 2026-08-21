package store

import (
	"context"
	"sync"
)

// PreviewSessionStore tracks request lifecycle state for preview work.
type PreviewSessionStore struct {
	mu       sync.Mutex
	firstCtx context.Context
}

func NewPreviewSessionStore() *PreviewSessionStore {
	return &PreviewSessionStore{}
}

func (s *PreviewSessionStore) SessionContext(ctx context.Context) context.Context {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.firstCtx == nil {
		s.firstCtx = ctx
	}
	return s.firstCtx
}
