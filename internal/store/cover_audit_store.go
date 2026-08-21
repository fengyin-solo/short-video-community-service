package store

import (
	"sync"

	"shortvideo/internal/model"
)

type CoverAuditStore struct {
	mu      sync.RWMutex
	entries map[string]*model.CoverRenderContext
}

func NewCoverAuditStore() *CoverAuditStore {
	return &CoverAuditStore{entries: make(map[string]*model.CoverRenderContext)}
}

func (s *CoverAuditStore) Save(requestID string, ctx *model.CoverRenderContext) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[requestID] = ctx
}

func (s *CoverAuditStore) Get(requestID string) (model.CoverRenderResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ctx, ok := s.entries[requestID]
	if !ok {
		return model.CoverRenderResult{}, false
	}
	return model.CoverRenderResult{RequestID: requestID, OwnerID: ctx.OwnerID, VideoID: ctx.VideoID}, true
}
