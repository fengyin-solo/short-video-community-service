package store

import "shortvideo/internal/model"

type PlaylistPublishAudit struct {
	entries []model.PlaylistPublication
}

func NewPlaylistPublishAudit() *PlaylistPublishAudit { return &PlaylistPublishAudit{} }

func (a *PlaylistPublishAudit) Append(entry model.PlaylistPublication) {
	a.entries = append(a.entries, entry)
}

func (a *PlaylistPublishAudit) Entries() []model.PlaylistPublication {
	return append([]model.PlaylistPublication(nil), a.entries...)
}
