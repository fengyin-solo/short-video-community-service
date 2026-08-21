package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"shortvideo/internal/coordination"
	"shortvideo/internal/model"
	"shortvideo/pkg/regiondelivery"
)

func TestVideoDeliveryBatchHandlesDuplicateAndLateRegions(t *testing.T) {
	runner := regiondelivery.NewRunner()
	svc := NewVideoDeliveryService(coordination.NewVideoDeliveryCoordinator(runner))
	ready := make(chan struct{})
	close(ready)
	results, err := svc.Deliver(context.Background(), []model.VideoDeliveryJob{
		{Region: "edge-a", Duplicate: true, Release: ready},
		{Region: "edge-b", Release: ready},
	}, 40*time.Millisecond)
	if err != nil {
		t.Errorf("completed delivery batch timed out: %v", err)
	}
	if len(results) != 1 || results[0].Region != "edge-b" {
		t.Errorf("completed region result was not preserved: %v", results)
	}

	slow := make(chan struct{})
	_, err = svc.Deliver(context.Background(), []model.VideoDeliveryJob{{Region: "edge-slow", Release: slow}}, 40*time.Millisecond)
	if !errors.Is(err, ErrVideoDeliveryTimeout) {
		t.Fatalf("slow region should time out, got %v", err)
	}
	close(slow)
	select {
	case panicText := <-runner.Panics():
		if strings.Contains(panicText, "send on closed channel") {
			t.Errorf("late delivery sent to a closed result channel: %s", panicText)
		}
	case <-time.After(80 * time.Millisecond):
	}
}
