package worker

import (
	"context"
	"errors"
	"fmt"

	"shortvideo/internal/model"
	"shortvideo/internal/scheduler"
	"shortvideo/pkg/safetyscan"
)

type SafetyScanWorker struct {
	client    *safetyscan.Client
	scheduler *scheduler.ScanRetryScheduler
	panics    chan string
}

func NewSafetyScanWorker(client *safetyscan.Client, scheduler *scheduler.ScanRetryScheduler) *SafetyScanWorker {
	return &SafetyScanWorker{client: client, scheduler: scheduler, panics: make(chan string, 2)}
}

func (w *SafetyScanWorker) Panics() <-chan string { return w.panics }

func (w *SafetyScanWorker) Run(ctx context.Context, videoID string, results chan<- model.SafetyScanResult) {
	go func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				w.panics <- fmt.Sprint(recovered)
			}
		}()
		result, err := w.client.Scan(ctx, videoID)
		if errors.Is(err, safetyscan.ErrTemporary) {
			w.scheduler.Schedule(ctx, func() { w.Run(ctx, videoID, results) })
			return
		}
		if err == nil {
			results <- result
		}
	}()
}
