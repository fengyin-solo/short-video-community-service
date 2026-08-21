package handler

import (
	"testing"

	"shortvideo/internal/service"
	"shortvideo/internal/store"
	"shortvideo/pkg/transcode"
)

func TestRetrySuccessSurvivesLateTranscodeCallback(t *testing.T) {
	w := &transcode.Worker{}
	h := NewTranscodeHandler(service.NewTranscodeService(store.NewTranscodeStore(), w))
	release, done := make(chan struct{}), make(chan struct{}, 1)
	h.Retry("video-render", release, done)
	close(release)
	<-done
	detail, listed := h.Detail("video-render"), h.Listed("video-render")
	if w.SideEffects != 1 {
		t.Errorf("retry executed transcode side effect %d times", w.SideEffects)
	}
	if detail.Status != "succeeded" || detail.Version != 2 {
		t.Errorf("late callback reverted successful transcode: %+v", detail)
	}
	if listed.Status != detail.Status || listed.Version != detail.Version {
		t.Errorf("list and detail disagree: list=%+v detail=%+v", listed, detail)
	}
}
