package store

import "sync"

type DeliveryStore struct {
	mu      sync.Mutex
	records map[string]int
}

func NewDeliveryStore() *DeliveryStore { return &DeliveryStore{records: map[string]int{}} }

func (s *DeliveryStore) Run(videoID string, send func() error) error {
	s.mu.Lock()
	s.records[videoID]++
	s.mu.Unlock()
	return send()
}

func (s *DeliveryStore) Records(videoID string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.records[videoID]
}
