package service

import (
	"reflect"
	"testing"

	"shortvideo/internal/store"
	"shortvideo/pkg/tagparse"
)

func TestTagBatchOwnershipAcrossAsyncIndexing(t *testing.T) {
	s := NewTagBatchService(tagparse.NewParser(), store.NewTagBatchCache())
	release := make(chan struct{})
	done := make(chan []string, 1)
	first := s.Submit("batch-first", []string{"travel", "night"}, release, done)
	s.Submit("batch-second", []string{"food", "quick"}, make(chan struct{}), make(chan []string, 1))
	close(release)
	indexed := <-done
	want := []string{"travel", "night"}
	if !reflect.DeepEqual(first, want) || !reflect.DeepEqual(indexed, want) {
		t.Errorf("first tag batch changed after second parse: response=%v indexed=%v", first, indexed)
	}
	if cached := s.Cached("batch-first"); !reflect.DeepEqual(cached, want) {
		t.Errorf("first tag cache was overwritten: %v", cached)
	}
}
