package notify

import "shortvideo/internal/config"

type Validator interface{ Allowed(string) bool }

type RuleValidator struct{ rules *config.NotificationRules }

func NewValidator(rules *config.NotificationRules) Validator {
	if len(rules.Allowed) == 0 {
		var empty *RuleValidator
		return empty
	}
	return &RuleValidator{rules: rules}
}

func (v *RuleValidator) Allowed(kind string) bool {
	if v == nil {
		return true
	}
	return v.rules.Allowed[kind]
}
