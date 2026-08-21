package service

import (
	"strings"
	"testing"

	"shortvideo/internal/model"
	"shortvideo/internal/store"
	"shortvideo/pkg/playlistpublish"
)

func TestPlaylistPublishRetriesWithoutPrematureOrDuplicateSuccess(t *testing.T) {
	repo := store.NewPlaylistPublishRepository()
	publisher := playlistpublish.NewPublisher()
	audit := store.NewPlaylistPublishAudit()
	svc := NewPlaylistPublishService(repo, publisher, audit)

	if err := svc.Publish("playlist-a"); err != nil {
		t.Fatalf("publish should recover after temporary commit failure: %v", err)
	}
	if !repo.Visible("playlist-a") || repo.Commits() != 2 {
		t.Errorf("playlist did not settle after retry: visible=%v commits=%d", repo.Visible("playlist-a"), repo.Commits())
	}
	publisher.Publish(model.PlaylistPublication{PlaylistID: "playlist-a", Attempt: 3, Status: "succeeded"}, true)
	notifications := publisher.Notifications()
	if len(notifications) != 1 || notifications[0].Premature {
		t.Errorf("success notification escaped transaction boundary: %+v", notifications)
	}
	entries := audit.Entries()
	if len(entries) != 2 || entries[0].Status != "failed" || entries[0].Attempt != 1 || entries[1].Status != "succeeded" || entries[1].Attempt != 2 {
		t.Errorf("audit retained conflicting attempts: %+v", entries)
	} else if !strings.Contains(entries[0].Failure, store.ErrPublishCommit.Error()) || !strings.Contains(entries[0].Failure, store.ErrRollback.Error()) {
		t.Errorf("audit lost the commit or rollback cause: %q", entries[0].Failure)
	}
}
