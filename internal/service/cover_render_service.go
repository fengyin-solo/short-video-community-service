package service

import (
	"shortvideo/internal/model"
	"shortvideo/internal/pool"
	"shortvideo/internal/store"
	"shortvideo/pkg/cover"
)

type CoverRenderService struct {
	contexts *pool.CoverContextPool
	renderer *cover.Renderer
	audits   *store.CoverAuditStore
}

func NewCoverRenderService(contexts *pool.CoverContextPool, renderer *cover.Renderer, audits *store.CoverAuditStore) *CoverRenderService {
	return &CoverRenderService{contexts: contexts, renderer: renderer, audits: audits}
}

func (s *CoverRenderService) Submit(requestID, ownerID, videoID string, release <-chan struct{}) <-chan model.CoverRenderResult {
	ctx := s.contexts.Acquire(requestID, ownerID, videoID)
	s.audits.Save(requestID, ctx)
	done := s.renderer.Render(ctx, release)
	s.contexts.Release(ctx)
	return done
}

func (s *CoverRenderService) Audit(requestID string) (model.CoverRenderResult, bool) {
	return s.audits.Get(requestID)
}
