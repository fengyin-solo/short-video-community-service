package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"shortvideo/internal/model"
	"shortvideo/internal/service"
	"shortvideo/internal/store"
	previewclient "shortvideo/pkg/preview"
)

func TestPreviewCancellationIsolation(t *testing.T) {
	client := previewclient.NewClient()
	h := NewPreviewHandler(service.NewPreviewService(store.NewPreviewSessionStore(), client))

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	startedAt := time.Now()
	go func() {
		_, err := h.Prepare(ctx, model.PreviewRequest{VideoID: "video-first", DelayMS: 180})
		done <- err
	}()
	<-client.Started
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("canceled preview returned %v, want context canceled", err)
		}
		if elapsed := time.Since(startedAt); elapsed > 100*time.Millisecond {
			t.Fatalf("canceled preview kept working for %s", elapsed)
		}
	case <-time.After(120 * time.Millisecond):
		t.Fatal("canceled preview did not stop promptly")
	}

	result, err := h.Prepare(context.Background(), model.PreviewRequest{VideoID: "video-next", DelayMS: 5})
	if err != nil {
		t.Fatalf("next preview inherited the previous cancellation: %v", err)
	}
	if !result.Ready || result.VideoID != "video-next" {
		t.Fatalf("unexpected next preview result: %+v", result)
	}
}
