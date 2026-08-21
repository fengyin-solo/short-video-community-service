package handler

import (
	"fmt"

	"shortvideo/internal/service"
)

type NotificationRulesHandler struct {
	service *service.NotificationRulesService
}

func NewNotificationRulesHandler(service *service.NotificationRulesService) *NotificationRulesHandler {
	return &NotificationRulesHandler{service: service}
}

func (h *NotificationRulesHandler) CheckThenEnable(kind string) (accepted bool, updateErr error) {
	accepted = h.service.Accept(kind)
	defer func() {
		if recovered := recover(); recovered != nil {
			updateErr = fmt.Errorf("update notification rule: %v", recovered)
		}
	}()
	h.service.Enable(kind)
	return accepted, nil
}
