package service

import (
	"context"
	"errors"
	"time"

	"shortvideo/internal/coordination"
	"shortvideo/internal/model"
)

var ErrVideoDeliveryTimeout = errors.New("video delivery timed out")

type VideoDeliveryService struct {
	coordinator *coordination.VideoDeliveryCoordinator
}

func NewVideoDeliveryService(coordinator *coordination.VideoDeliveryCoordinator) *VideoDeliveryService {
	return &VideoDeliveryService{coordinator: coordinator}
}

func (s *VideoDeliveryService) Deliver(ctx context.Context, jobs []model.VideoDeliveryJob, timeout time.Duration) ([]model.VideoDeliveryResult, error) {
	batch := s.coordinator.Start(ctx, jobs)
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	results := make([]model.VideoDeliveryResult, 0, len(jobs))
	for {
		select {
		case result, ok := <-batch.Results:
			if !ok {
				return results, nil
			}
			results = append(results, result)
		case <-timer.C:
			batch.Abort()
			return results, ErrVideoDeliveryTimeout
		case <-ctx.Done():
			batch.Abort()
			return results, ctx.Err()
		}
	}
}
