package dispatch

import (
	"context"

	"shortvideo/internal/model"
	"shortvideo/internal/worker"
)

type SafetyScanDispatcher struct{ worker *worker.SafetyScanWorker }

func NewSafetyScanDispatcher(worker *worker.SafetyScanWorker) *SafetyScanDispatcher {
	return &SafetyScanDispatcher{worker: worker}
}

func (d *SafetyScanDispatcher) Dispatch(_ context.Context, videoID string, results chan<- model.SafetyScanResult) {
	d.worker.Run(context.Background(), videoID, results)
}
