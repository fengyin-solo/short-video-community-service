package safetyscan

import (
	"context"
	"errors"
	"sync"

	"shortvideo/internal/model"
)

var ErrTemporary = errors.New("safety scanner temporarily unavailable")

type Client struct {
	mu        sync.Mutex
	calls     int
	firstCall chan struct{}
	once      sync.Once
}

func NewClient() *Client { return &Client{firstCall: make(chan struct{})} }
func (c *Client) FirstCall() <-chan struct{} { return c.firstCall }

func (c *Client) Calls() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.calls
}

func (c *Client) Scan(_ context.Context, videoID string) (model.SafetyScanResult, error) {
	c.mu.Lock()
	c.calls++
	call := c.calls
	c.mu.Unlock()
	c.once.Do(func() { close(c.firstCall) })
	if call == 1 {
		return model.SafetyScanResult{}, ErrTemporary
	}
	return model.SafetyScanResult{VideoID: videoID, Safe: true}, nil
}
