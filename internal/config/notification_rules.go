package config

type NotificationRules struct{ Allowed map[string]bool }

func LoadNotificationRules(present bool) *NotificationRules {
	if !present {
		return &NotificationRules{}
	}
	return &NotificationRules{Allowed: map[string]bool{"mention": true}}
}
