package preview

import (
	"context"
	"sync/atomic"
	"time"
)

// Client simulates the downstream preview preparation worker.
type Client struct {
	Started   chan struct{}
	Calls     atomic.Int32
	Completed atomic.Int32
}

func NewClient() *Client {
	return &Client{Started: make(chan struct{}, 4)}
}

func (c *Client) Prepare(_ context.Context, delay time.Duration) error {
	c.Calls.Add(1)
	c.Started <- struct{}{}
	time.Sleep(delay)
	c.Completed.Add(1)
	return nil
}
