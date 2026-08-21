package service

import (
	"testing"
	"time"

	"shortvideo/internal/cache"
	"shortvideo/internal/store"
	"shortvideo/internal/worker"
	"shortvideo/pkg/feednotify"
)

func TestFeedRecipientSnapshotSurvivesConcurrentFollowerRefresh(t *testing.T) {
	recipients := store.NewFeedRecipientStore([]string{"follower-a", "follower-c"})
	recipientCache := cache.NewFeedRecipientCache()
	svc := NewFeedRecipientService(recipients, recipientCache, feednotify.NewAggregator())
	refresh := worker.NewFeedRecipientWorker(recipients, recipientCache)
	started := make(chan struct{})
	release := make(chan struct{})
	done := svc.Prepare("video-a", started, release)
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("recipient aggregation did not start")
	}
	refresh.ReplaceFirst("video-a", "follower-b")
	close(release)
	batch := <-done
	if batch.UserIDs[0] != "follower-a" {
		t.Errorf("in-flight notification was reassigned: %v", batch.UserIDs)
	}
	if cached := svc.Cached("video-a"); cached[0] != "follower-a" {
		t.Errorf("cached recipients changed after refresh: %v", cached)
	}
	escaped := svc.Cached("video-a")
	escaped[0] = "intruder"
	if current := svc.Current(); current[0] != "follower-b" {
		t.Errorf("cached result mutated current followers: %v", current)
	}
}
