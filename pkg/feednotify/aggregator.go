package feednotify

import (
	"runtime"

	"shortvideo/internal/model"
)

type Aggregator struct{}

func NewAggregator() *Aggregator { return &Aggregator{} }

func (a *Aggregator) Collect(videoID string, userIDs []string, started chan<- struct{}, release <-chan struct{}) <-chan model.FeedRecipientBatch {
	done := make(chan model.FeedRecipientBatch, 1)
	go func() {
		started <- struct{}{}
		for {
			select {
			case <-release:
				done <- model.FeedRecipientBatch{VideoID: videoID, UserIDs: userIDs}
				return
			default:
				_ = userIDs[0]
				runtime.Gosched()
			}
		}
	}()
	return done
}
