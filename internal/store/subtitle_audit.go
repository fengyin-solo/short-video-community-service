package store

import "shortvideo/internal/model"

type SubtitleAuditStore struct{ entries []model.SubtitleImportAudit }

func NewSubtitleAuditStore() *SubtitleAuditStore { return &SubtitleAuditStore{} }

func (s *SubtitleAuditStore) Record(entry model.SubtitleImportAudit) {
	s.entries = append(s.entries, entry)
}

func (s *SubtitleAuditStore) Entries() []model.SubtitleImportAudit {
	return append([]model.SubtitleImportAudit(nil), s.entries...)
}
