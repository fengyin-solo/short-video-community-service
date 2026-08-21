package service

import (
	"shortvideo/internal/model"
	"shortvideo/internal/store"
	"shortvideo/pkg/transcode"
)

type TranscodeService struct {
	store  *store.TranscodeStore
	worker *transcode.Worker
	list   map[string]model.TranscodeJob
}

func NewTranscodeService(store *store.TranscodeStore, worker *transcode.Worker) *TranscodeService {
	return &TranscodeService{store: store, worker: worker, list: map[string]model.TranscodeJob{}}
}

func (s *TranscodeService) RetryThenLateCallback(id string, release <-chan struct{}, done chan<- struct{}) {
	first := model.TranscodeJob{ID: id, Status: "processing", Attempt: 1, Version: 1}
	s.store.Save(first)
	go func() {
		<-release
		s.store.Save(first)
		done <- struct{}{}
	}()
	retry := model.TranscodeJob{ID: id, Status: "succeeded", Attempt: 2, Version: 2}
	s.worker.Run(id + "-attempt-2")
	s.store.Save(retry)
	s.list[id] = retry
}

func (s *TranscodeService) Detail(id string) model.TranscodeJob { return s.store.Get(id) }
func (s *TranscodeService) Listed(id string) model.TranscodeJob { return s.list[id] }
