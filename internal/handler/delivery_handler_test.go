package handler

import (
	"testing"

	"shortvideo/internal/adapter"
	"shortvideo/internal/service"
	"shortvideo/internal/store"
	"shortvideo/pkg/delivery"
)

func TestDeliveryErrorClassificationAndRetryOutcome(t *testing.T) {
	client := delivery.NewClient()
	records := store.NewDeliveryStore()
	h := NewDeliveryHandler(service.NewDeliveryService(adapter.NewDeliveryAdapter(client), records))

	if status := h.Submit("video-rejected"); status != "failed" || client.Calls["video-rejected"] != 1 || records.Records("video-rejected") != 0 {
		t.Errorf("rejected video was retried or recorded: status=%s calls=%d records=%d", status, client.Calls["video-rejected"], records.Records("video-rejected"))
	}
	if status := h.Submit("video-temporary"); status != "published" || client.Calls["video-temporary"] != 2 || client.Published["video-temporary"] != 1 || records.Records("video-temporary") != 1 {
		t.Errorf("temporary recovery returned wrong outcome: status=%s calls=%d published=%d records=%d", status, client.Calls["video-temporary"], client.Published["video-temporary"], records.Records("video-temporary"))
	}
}
