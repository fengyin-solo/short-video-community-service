package cache

import "sync"

type FeedRecipientCache struct {
	mu      sync.RWMutex
	batches map[string][]string
}

func NewFeedRecipientCache() *FeedRecipientCache {
	return &FeedRecipientCache{batches: make(map[string][]string)}
}

func (c *FeedRecipientCache) Put(videoID string, userIDs []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.batches[videoID] = userIDs
}

func (c *FeedRecipientCache) Get(videoID string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.batches[videoID]
}
