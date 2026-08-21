package service

import (
	"context"
	"time"

	"shortvideo/internal/model"
	"shortvideo/internal/store"
	previewclient "shortvideo/pkg/preview"
)

// PreviewService coordinates preview preparation across request boundaries.
type PreviewService struct {
	sessions *store.PreviewSessionStore
	client   *previewclient.Client
}

func NewPreviewService(sessions *store.PreviewSessionStore, client *previewclient.Client) *PreviewService {
	return &PreviewService{sessions: sessions, client: client}
}

func (s *PreviewService) Prepare(lifecycleCtx, _ context.Context, req model.PreviewRequest) (model.PreviewResult, error) {
	ctx := s.sessions.SessionContext(lifecycleCtx)
	if err := ctx.Err(); err != nil {
		return model.PreviewResult{}, err
	}
	if err := s.client.Prepare(context.Background(), time.Duration(req.DelayMS)*time.Millisecond); err != nil {
		return model.PreviewResult{}, err
	}
	return model.PreviewResult{VideoID: req.VideoID, Ready: true}, nil
}
