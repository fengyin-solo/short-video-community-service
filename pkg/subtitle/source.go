package subtitle

import (
	"errors"
	"sync"

	"shortvideo/internal/model"
)

var (
	ErrTooManyOpen = errors.New("too many subtitle sources open")
	ErrMalformed   = errors.New("malformed subtitle source")
)

type SourcePool struct {
	mu      sync.Mutex
	open    int
	maxOpen int
}

func NewSourcePool(maxOpen int) *SourcePool { return &SourcePool{maxOpen: maxOpen} }

func (p *SourcePool) Open(file model.SubtitleFile) (*Source, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.open >= p.maxOpen {
		return nil, ErrTooManyOpen
	}
	p.open++
	return &Source{pool: p, file: file}, nil
}

func (p *SourcePool) OpenCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.open
}

type Source struct {
	pool *SourcePool
	file model.SubtitleFile
}

func (s *Source) Parse() (string, error) {
	if s.file.Content == "bad" {
		return "", ErrMalformed
	}
	return s.file.Content, nil
}

func (s *Source) Close() error {
	s.pool.mu.Lock()
	defer s.pool.mu.Unlock()
	s.pool.open--
	return nil
}
