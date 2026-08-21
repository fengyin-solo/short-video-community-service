package store

type TagBatchCache struct{ batches map[string][]string }

func NewTagBatchCache() *TagBatchCache                { return &TagBatchCache{batches: map[string][]string{}} }
func (c *TagBatchCache) Put(id string, tags []string) { c.batches[id] = tags }
func (c *TagBatchCache) Get(id string) []string       { return c.batches[id] }
