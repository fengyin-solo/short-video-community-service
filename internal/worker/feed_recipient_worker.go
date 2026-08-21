package worker

import (
	"shortvideo/internal/cache"
	"shortvideo/internal/store"
)

type FeedRecipientWorker struct {
	store *store.FeedRecipientStore
	cache *cache.FeedRecipientCache
}

func NewFeedRecipientWorker(store *store.FeedRecipientStore, cache *cache.FeedRecipientCache) *FeedRecipientWorker {
	return &FeedRecipientWorker{store: store, cache: cache}
}

func (w *FeedRecipientWorker) ReplaceFirst(videoID, userID string) {
	w.store.ReplaceFirst(userID)
	cached := w.cache.Get(videoID)
	if len(cached) > 0 {
		cached[0] = userID
	}
}
