package store

import (
	"sync"

	"shortvideo/internal/model"
)

type TranscodeStore struct {
	mu   sync.Mutex
	jobs map[string]model.TranscodeJob
}

func NewTranscodeStore() *TranscodeStore {
	return &TranscodeStore{jobs: map[string]model.TranscodeJob{}}
}
func (s *TranscodeStore) Save(job model.TranscodeJob) {
	s.mu.Lock()
	s.jobs[job.ID] = job
	s.mu.Unlock()
}
func (s *TranscodeStore) Get(id string) model.TranscodeJob {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.jobs[id]
}
