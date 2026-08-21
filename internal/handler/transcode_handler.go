package handler

import (
	"shortvideo/internal/model"
	"shortvideo/internal/service"
)

type TranscodeHandler struct{ service *service.TranscodeService }

func NewTranscodeHandler(service *service.TranscodeService) *TranscodeHandler {
	return &TranscodeHandler{service: service}
}
func (h *TranscodeHandler) Retry(id string, release <-chan struct{}, done chan<- struct{}) {
	h.service.RetryThenLateCallback(id, release, done)
}
func (h *TranscodeHandler) Detail(id string) model.TranscodeJob { return h.service.Detail(id) }
func (h *TranscodeHandler) Listed(id string) model.TranscodeJob { return h.service.Listed(id) }
