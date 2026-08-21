package service

import (
	"shortvideo/internal/store"
	"shortvideo/pkg/tagindex"
	"shortvideo/pkg/tagparse"
)

type TagBatchService struct {
	parser *tagparse.Parser
	cache  *store.TagBatchCache
}

func NewTagBatchService(parser *tagparse.Parser, cache *store.TagBatchCache) *TagBatchService {
	return &TagBatchService{parser: parser, cache: cache}
}

func (s *TagBatchService) Submit(id string, values []string, release <-chan struct{}, done chan<- []string) []string {
	tags := s.parser.Parse(values)
	s.cache.Put(id, tags)
	go tagindex.Job{Tags: tags, Release: release, Done: done}.Run()
	return tags
}

func (s *TagBatchService) Cached(id string) []string { return s.cache.Get(id) }
