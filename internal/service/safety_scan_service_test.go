package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"shortvideo/internal/dispatch"
	"shortvideo/internal/scheduler"
	"shortvideo/internal/worker"
	"shortvideo/pkg/safetyscan"
)

func TestSafetyScanStopsQueuedRetryAfterCancel(t *testing.T) {
	gate := make(chan struct{})
	client := safetyscan.NewClient()
	retries := scheduler.NewScanRetryScheduler(gate)
	scanWorker := worker.NewSafetyScanWorker(client, retries)
	svc := NewSafetyScanService(dispatch.NewSafetyScanDispatcher(scanWorker))
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := svc.Scan(ctx, "video-a")
		done <- err
	}()
	select {
	case <-client.FirstCall():
	case <-time.After(time.Second):
		t.Fatal("initial safety scan did not start")
	}
	select {
	case <-retries.Scheduled():
	case <-time.After(time.Second):
		t.Fatal("temporary failure did not schedule a retry")
	}
	cancel()
	if err := <-done; !errors.Is(err, context.Canceled) {
		t.Fatalf("scan should return cancellation, got %v", err)
	}
	close(gate)
	deadline := time.After(time.Second)
	for client.Calls() == 1 {
		select {
		case <-deadline:
			t.Fatal("retry did not run after cancellation")
		default:
			time.Sleep(time.Millisecond)
		}
	}
	if calls := client.Calls(); calls != 1 {
		t.Errorf("scanner calls grew after request returned: %d", calls)
	}
	select {
	case panicText := <-scanWorker.Panics():
		if strings.Contains(panicText, "send on closed channel") {
			t.Errorf("late scan result hit a closed channel: %s", panicText)
		}
	case <-time.After(100 * time.Millisecond):
	}
}
