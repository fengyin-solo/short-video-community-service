package handler

import "shortvideo/internal/service"

type DeliveryHandler struct{ service *service.DeliveryService }

func NewDeliveryHandler(service *service.DeliveryService) *DeliveryHandler {
	return &DeliveryHandler{service: service}
}

func (h *DeliveryHandler) Submit(videoID string) string {
	if err := h.service.Deliver(videoID); err != nil {
		return "failed"
	}
	return "published"
}
