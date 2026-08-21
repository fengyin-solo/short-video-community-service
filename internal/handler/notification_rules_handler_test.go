package handler

import (
	"testing"

	"shortvideo/internal/config"
	"shortvideo/internal/service"
)

func TestMissingNotificationRulesRejectAndInitialize(t *testing.T) {
	h := NewNotificationRulesHandler(service.NewNotificationRulesService(config.LoadNotificationRules(false)))
	accepted, err := h.CheckThenEnable("comment")
	if accepted {
		t.Error("missing rules silently accepted comment notification")
	}
	if err != nil {
		t.Errorf("enabling default rule panicked: %v", err)
	}
}
