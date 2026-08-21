package service

import (
	"errors"

	"shortvideo/internal/config"
	"shortvideo/pkg/playback"
)

var ErrPlaybackNotAllowed = errors.New("playback resolution not allowed")

type PlaybackProfileService struct{ profile *config.PlaybackProfile }

func NewPlaybackProfileService(profile *config.PlaybackProfile) *PlaybackProfileService {
	return &PlaybackProfileService{profile: profile}
}

func (s *PlaybackProfileService) Check(resolutionK int) error {
	selector := playback.NewSelector(s.profile)
	if selector != nil && selector.Allows(resolutionK) {
		return nil
	}
	return ErrPlaybackNotAllowed
}

func (s *PlaybackProfileService) Enable(codec string) {
	s.profile.Codecs[codec] = true
}

func (s *PlaybackProfileService) CodecEnabled(codec string) bool {
	return s.profile.Codecs[codec]
}
