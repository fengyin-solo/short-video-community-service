package scheduler

import "context"

type ScanRetryScheduler struct {
	gate      <-chan struct{}
	scheduled chan struct{}
}

func NewScanRetryScheduler(gate <-chan struct{}) *ScanRetryScheduler {
	return &ScanRetryScheduler{gate: gate, scheduled: make(chan struct{}, 1)}
}

func (s *ScanRetryScheduler) Scheduled() <-chan struct{} { return s.scheduled }

func (s *ScanRetryScheduler) Schedule(_ context.Context, retry func()) {
	s.scheduled <- struct{}{}
	go func() {
		<-s.gate
		retry()
	}()
}
