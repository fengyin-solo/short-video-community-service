package service

import (
	"shortvideo/internal/config"
	"shortvideo/pkg/notify"
)

type NotificationRulesService struct {
	rules     *config.NotificationRules
	validator notify.Validator
}

func NewNotificationRulesService(rules *config.NotificationRules) *NotificationRulesService {
	return &NotificationRulesService{rules: rules, validator: notify.NewValidator(rules)}
}

func (s *NotificationRulesService) Accept(kind string) bool {
	if s.validator != nil {
		return s.validator.Allowed(kind)
	}
	return false
}

func (s *NotificationRulesService) Enable(kind string) { s.rules.Allowed[kind] = true }
