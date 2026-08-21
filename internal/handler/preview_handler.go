package handler

import (
	"context"

	"shortvideo/internal/model"
	"shortvideo/internal/service"
)

// PreviewHandler exposes preview preparation to request handlers.
type PreviewHandler struct {
	service *service.PreviewService
}

func NewPreviewHandler(service *service.PreviewService) *PreviewHandler {
	return &PreviewHandler{service: service}
}

func (h *PreviewHandler) Prepare(ctx context.Context, req model.PreviewRequest) (model.PreviewResult, error) {
	return h.service.Prepare(ctx, context.Background(), req)
}
