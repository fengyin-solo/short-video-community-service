package service

import (
	"shortvideo/internal/cache"
	"shortvideo/internal/model"
	"shortvideo/internal/store"
	"shortvideo/pkg/feednotify"
)

type FeedRecipientService struct {
	store      *store.FeedRecipientStore
	cache      *cache.FeedRecipientCache
	aggregator *feednotify.Aggregator
}

func NewFeedRecipientService(store *store.FeedRecipientStore, cache *cache.FeedRecipientCache, aggregator *feednotify.Aggregator) *FeedRecipientService {
	return &FeedRecipientService{store: store, cache: cache, aggregator: aggregator}
}

func (s *FeedRecipientService) Prepare(videoID string, started chan<- struct{}, release <-chan struct{}) <-chan model.FeedRecipientBatch {
	userIDs := s.store.Snapshot()
	s.cache.Put(videoID, userIDs)
	return s.aggregator.Collect(videoID, userIDs, started, release)
}

func (s *FeedRecipientService) Cached(videoID string) []string { return s.cache.Get(videoID) }
func (s *FeedRecipientService) Current() []string              { return s.store.Snapshot() }
