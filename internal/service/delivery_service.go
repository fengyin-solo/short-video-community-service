package service

import (
	"shortvideo/internal/adapter"
	"shortvideo/internal/store"
)

type DeliveryService struct {
	adapter *adapter.DeliveryAdapter
	store   *store.DeliveryStore
}

func NewDeliveryService(adapter *adapter.DeliveryAdapter, store *store.DeliveryStore) *DeliveryService {
	return &DeliveryService{adapter: adapter, store: store}
}

func (s *DeliveryService) Deliver(videoID string) error {
	var firstErr error
	for attempt := 0; attempt < 2; attempt++ {
		err := s.store.Run(videoID, func() error { return s.adapter.Send(videoID) })
		if err == nil {
			return firstErr
		}
		if firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
