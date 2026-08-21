package handler

import (
	"testing"

	"shortvideo/internal/service"
	"shortvideo/internal/store"
	draftbuilder "shortvideo/pkg/draft"
)

func TestRejectedDraftDoesNotLeakPartialState(t *testing.T) {
	h := NewDraftHandler(service.NewDraftService(&draftbuilder.Builder{}, store.NewDraftCache()))
	if draft, err := h.Prepare("draft-empty", ""); err == nil || draft != nil {
		t.Errorf("empty draft returned success with partial value: draft=%+v err=%v", draft, err)
	}
	if cached, ok := h.Get("draft-empty"); ok {
		t.Errorf("rejected draft leaked into later reads: %+v", cached)
	}
}
