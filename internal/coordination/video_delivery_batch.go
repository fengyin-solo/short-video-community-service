package coordination

import (
	"context"
	"sync"

	"shortvideo/internal/model"
	"shortvideo/pkg/regiondelivery"
)

type VideoDeliveryBatch struct {
	Results <-chan model.VideoDeliveryResult
	output  chan model.VideoDeliveryResult
	once    sync.Once
}

type VideoDeliveryCoordinator struct {
	runner *regiondelivery.Runner
}

func NewVideoDeliveryCoordinator(runner *regiondelivery.Runner) *VideoDeliveryCoordinator {
	return &VideoDeliveryCoordinator{runner: runner}
}

func (c *VideoDeliveryCoordinator) Start(ctx context.Context, jobs []model.VideoDeliveryJob) *VideoDeliveryBatch {
	output := make(chan model.VideoDeliveryResult, len(jobs))
	var wg sync.WaitGroup
	wg.Add(len(jobs))
	for _, job := range jobs {
		c.runner.Run(ctx, job, output, &wg)
	}
	go func() {
		wg.Wait()
		close(output)
	}()
	return &VideoDeliveryBatch{Results: output, output: output}
}

func (b *VideoDeliveryBatch) Abort() {
	b.once.Do(func() { close(b.output) })
}
