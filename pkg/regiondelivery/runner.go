package regiondelivery

import (
	"context"
	"fmt"
	"sync"

	"shortvideo/internal/model"
)

type Runner struct {
	panics chan string
}

func NewRunner() *Runner { return &Runner{panics: make(chan string, 4)} }

func (r *Runner) Panics() <-chan string { return r.panics }

func (r *Runner) Run(_ context.Context, job model.VideoDeliveryJob, results chan<- model.VideoDeliveryResult, wg *sync.WaitGroup) {
	go func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				r.panics <- fmt.Sprint(recovered)
			}
		}()
		if job.Duplicate {
			return
		}
		<-job.Release
		results <- model.VideoDeliveryResult{Region: job.Region}
		wg.Done()
	}()
}
