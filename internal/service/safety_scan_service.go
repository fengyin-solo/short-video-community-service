package service

import (
	"context"

	"shortvideo/internal/dispatch"
	"shortvideo/internal/model"
)

type SafetyScanService struct{ dispatcher *dispatch.SafetyScanDispatcher }

func NewSafetyScanService(dispatcher *dispatch.SafetyScanDispatcher) *SafetyScanService {
	return &SafetyScanService{dispatcher: dispatcher}
}

func (s *SafetyScanService) Scan(ctx context.Context, videoID string) (model.SafetyScanResult, error) {
	results := make(chan model.SafetyScanResult)
	s.dispatcher.Dispatch(ctx, videoID, results)
	select {
	case result := <-results:
		close(results)
		return result, nil
	case <-ctx.Done():
		close(results)
		return model.SafetyScanResult{}, ctx.Err()
	}
}
