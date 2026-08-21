package service_test

import (
	"errors"
	"testing"

	"shortvideo/internal/config"
	"shortvideo/internal/handler"
	"shortvideo/internal/service"
)

func TestUnknownPlaybackProfileRejects4KAndEnablesFallbackCodec(t *testing.T) {
	profile := config.LoadPlaybackProfile("")
	svc := service.NewPlaybackProfileService(profile)
	if err := svc.Check(4); !errors.Is(err, service.ErrPlaybackNotAllowed) {
		t.Errorf("unknown device was allowed 4K playback: %v", err)
	}
	h := handler.NewPlaybackProfileHandler(svc)
	if err := h.CheckThenEnable(2, "av1"); err != nil {
		t.Errorf("fallback codec update returned an internal error: %v", err)
	}
	if !svc.CodecEnabled("av1") {
		t.Error("fallback AV1 codec was not enabled")
	}
}
