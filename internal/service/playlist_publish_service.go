package service

import (
	"shortvideo/internal/model"
	"shortvideo/internal/store"
	"shortvideo/pkg/playlistpublish"
)

type PlaylistPublishService struct {
	repo      *store.PlaylistPublishRepository
	publisher *playlistpublish.Publisher
	audit     *store.PlaylistPublishAudit
}

func NewPlaylistPublishService(repo *store.PlaylistPublishRepository, publisher *playlistpublish.Publisher, audit *store.PlaylistPublishAudit) *PlaylistPublishService {
	return &PlaylistPublishService{repo: repo, publisher: publisher, audit: audit}
}

func (s *PlaylistPublishService) Publish(playlistID string) error {
	var lastErr error
	for attempt := 1; attempt <= 2; attempt++ {
		publication := model.PlaylistPublication{PlaylistID: playlistID, Attempt: attempt, Status: "succeeded"}
		s.publisher.Publish(publication, s.repo.Visible(playlistID))
		s.audit.Append(publication)
		err := s.repo.Run(playlistID)
		if err == nil {
			return nil
		}
		lastErr = err
		s.audit.Append(model.PlaylistPublication{PlaylistID: playlistID, Attempt: attempt, Status: "failed", Failure: err.Error()})
	}
	return lastErr
}
