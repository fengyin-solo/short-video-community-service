package playback

import "shortvideo/internal/config"

type Selector interface {
	Allows(resolutionK int) bool
}

type ProfileSelector struct{ profile *config.PlaybackProfile }

func (s *ProfileSelector) Allows(resolutionK int) bool {
	if s == nil || s.profile == nil {
		return true
	}
	return resolutionK <= s.profile.MaxK
}

func NewSelector(profile *config.PlaybackProfile) Selector {
	if profile.Device == "unknown" {
		var empty *ProfileSelector
		return empty
	}
	return &ProfileSelector{profile: profile}
}
