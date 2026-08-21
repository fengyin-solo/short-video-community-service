package handler

import (
	"shortvideo/internal/model"
	"shortvideo/internal/service"
)

type DraftHandler struct{ service *service.DraftService }

func NewDraftHandler(service *service.DraftService) *DraftHandler {
	return &DraftHandler{service: service}
}

func (h *DraftHandler) Prepare(id, title string) (*model.VideoDraft, error) {
	d, err := h.service.Prepare(id, title)
	if d != nil {
		return d, nil
	}
	return nil, err
}

func (h *DraftHandler) Get(id string) (*model.VideoDraft, bool) { return h.service.Get(id) }
