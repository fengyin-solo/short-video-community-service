package service

import (
	"errors"
	"testing"

	"shortvideo/internal/model"
	"shortvideo/internal/store"
	"shortvideo/pkg/subtitle"
)

func TestSubtitleBatchReleasesSourcesAndRollsBackMalformedImport(t *testing.T) {
	pool := subtitle.NewSourcePool(2)
	svc := NewSubtitleImportService(pool, store.NewSubtitleTransaction(), store.NewSubtitleAuditStore())
	valid := []model.SubtitleFile{{Name: "a", Content: "cue-a"}, {Name: "b", Content: "cue-b"}, {Name: "c", Content: "cue-c"}}
	if err := svc.Validate(valid); err != nil {
		t.Errorf("valid subtitle validation exhausted sources: %v", err)
	}
	if pool.OpenCount() != 0 {
		t.Errorf("subtitle sources remained open: %d", pool.OpenCount())
	}

	broken := []model.SubtitleFile{{Name: "first", Content: "cue-first"}, {Name: "broken", Content: "bad"}}
	if err := svc.Import("batch-a", broken); !errors.Is(err, subtitle.ErrMalformed) {
		t.Errorf("malformed subtitle error was lost: %v", err)
	}
	if got := svc.Visible(); len(got) != 0 {
		t.Errorf("partial subtitles became visible: %v", got)
	}
	audits := svc.Audits()
	if len(audits) != 1 || audits[0].Status != "failed" || audits[0].Failure != subtitle.ErrMalformed.Error() {
		t.Errorf("subtitle import audit reported wrong outcome: %+v", audits)
	}
}
