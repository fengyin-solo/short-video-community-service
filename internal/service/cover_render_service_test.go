package service

import (
	"testing"
	"time"

	"shortvideo/internal/model"
	"shortvideo/internal/pool"
	"shortvideo/internal/store"
	"shortvideo/pkg/cover"
)

func TestCoverRenderRequestsKeepOwnerAndVideoIsolated(t *testing.T) {
	started := make(chan string, 2)
	svc := NewCoverRenderService(pool.NewCoverContextPool(), cover.NewRenderer(started), store.NewCoverAuditStore())

	releaseA := make(chan struct{})
	doneA := svc.Submit("request-a", "owner-a", "video-a", releaseA)
	waitCoverStart(t, started, "request-a")

	doneB := svc.Submit("request-b", "owner-b", "video-b", nil)
	waitCoverStart(t, started, "request-b")
	resultB := waitCoverResult(t, doneB)
	close(releaseA)
	resultA := waitCoverResult(t, doneA)

	if resultA.RequestID != "request-a" || resultA.OwnerID != "owner-a" || resultA.VideoID != "video-a" {
		t.Errorf("first render used another request: %+v", resultA)
	}
	if resultB.RequestID != "request-b" || resultB.OwnerID != "owner-b" || resultB.VideoID != "video-b" {
		t.Errorf("second render mismatch: %+v", resultB)
	}
	auditA, ok := svc.Audit("request-a")
	if !ok || auditA.OwnerID != "owner-a" || auditA.VideoID != "video-a" {
		t.Errorf("first audit was reassigned: found=%v entry=%+v", ok, auditA)
	}
}

func waitCoverStart(t *testing.T, started <-chan string, want string) {
	t.Helper()
	select {
	case got := <-started:
		if got != want {
			t.Fatalf("render start: got %q want %q", got, want)
		}
	case <-time.After(time.Second):
		t.Fatal("render did not start")
	}
}

func waitCoverResult(t *testing.T, done <-chan model.CoverRenderResult) model.CoverRenderResult {
	t.Helper()
	select {
	case result := <-done:
		return result
	case <-time.After(time.Second):
		t.Fatal("render did not finish")
		return model.CoverRenderResult{}
	}
}
