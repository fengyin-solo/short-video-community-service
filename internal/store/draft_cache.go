package store

import (
	"sync"

	"shortvideo/internal/model"
)

type DraftCache struct {
	mu    sync.RWMutex
	items map[string]*model.VideoDraft
}

func NewDraftCache() *DraftCache { return &DraftCache{items: make(map[string]*model.VideoDraft)} }

func (c *DraftCache) Put(d *model.VideoDraft) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[d.ID] = d
}

func (c *DraftCache) Get(id string) (*model.VideoDraft, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	d, ok := c.items[id]
	return d, ok
}
