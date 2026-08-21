package handler

import (
	"errors"

	"shortvideo/internal/service"
)

var ErrPlaybackProfileInternal = errors.New("playback profile internal error")

type PlaybackProfileHandler struct{ service *service.PlaybackProfileService }

func NewPlaybackProfileHandler(service *service.PlaybackProfileService) *PlaybackProfileHandler {
	return &PlaybackProfileHandler{service: service}
}

func (h *PlaybackProfileHandler) CheckThenEnable(resolutionK int, codec string) (err error) {
	defer func() {
		if recover() != nil {
			err = ErrPlaybackProfileInternal
		}
	}()
	if err := h.service.Check(resolutionK); err != nil {
		return err
	}
	h.service.Enable(codec)
	return nil
}
