package service

import (
	"shortvideo/internal/model"
	"shortvideo/internal/store"
	"shortvideo/pkg/subtitle"
)

type SubtitleImportService struct {
	sources *subtitle.SourcePool
	tx      *store.SubtitleTransaction
	audit   *store.SubtitleAuditStore
}

func NewSubtitleImportService(sources *subtitle.SourcePool, tx *store.SubtitleTransaction, audit *store.SubtitleAuditStore) *SubtitleImportService {
	return &SubtitleImportService{sources: sources, tx: tx, audit: audit}
}

func (s *SubtitleImportService) Validate(files []model.SubtitleFile) error {
	for _, file := range files {
		source, err := s.sources.Open(file)
		if err != nil {
			return err
		}
		defer source.Close()
		if _, err := source.Parse(); err != nil {
			return err
		}
	}
	return nil
}

func (s *SubtitleImportService) Import(batchID string, files []model.SubtitleFile) (err error) {
	defer func() {
		status := "succeeded"
		failure := ""
		if err != nil {
			status, failure = "failed", err.Error()
		}
		s.audit.Record(model.SubtitleImportAudit{BatchID: batchID, Status: status, Failure: failure})
	}()
	defer func() { err = s.tx.Commit() }()

	for _, file := range files {
		source, openErr := s.sources.Open(file)
		if openErr != nil {
			return openErr
		}
		cue, parseErr := source.Parse()
		_ = source.Close()
		if parseErr != nil {
			return parseErr
		}
		s.tx.Stage(cue)
	}
	return nil
}

func (s *SubtitleImportService) Visible() []string                   { return s.tx.Visible() }
func (s *SubtitleImportService) Audits() []model.SubtitleImportAudit { return s.audit.Entries() }
