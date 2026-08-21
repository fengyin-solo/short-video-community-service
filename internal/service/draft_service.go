package service

import (
	"shortvideo/internal/model"
	"shortvideo/internal/store"
	draftbuilder "shortvideo/pkg/draft"
)

type DraftService struct {
	builder *draftbuilder.Builder
	cache   *store.DraftCache
}

func NewDraftService(builder *draftbuilder.Builder, cache *store.DraftCache) *DraftService {
	return &DraftService{builder: builder, cache: cache}
}

func (s *DraftService) Prepare(id, title string) (*model.VideoDraft, error) {
	d, err := s.builder.Build(id, title)
	if d != nil {
		s.cache.Put(d)
	}
	return d, err
}

func (s *DraftService) Get(id string) (*model.VideoDraft, bool) { return s.cache.Get(id) }
