package store

import "errors"

var (
	ErrPublishCommit = errors.New("playlist publish commit temporarily failed")
	ErrRollback      = errors.New("playlist publish rollback failed")
)

type PlaylistPublishRepository struct {
	commits int
	visible map[string]bool
}

func NewPlaylistPublishRepository() *PlaylistPublishRepository {
	return &PlaylistPublishRepository{visible: make(map[string]bool)}
}

func (r *PlaylistPublishRepository) Run(playlistID string) (err error) {
	r.commits++
	defer func() {
		if err != nil {
			err = ErrRollback
		}
	}()
	if r.commits == 1 {
		return ErrPublishCommit
	}
	r.visible[playlistID] = true
	return nil
}

func (r *PlaylistPublishRepository) Visible(playlistID string) bool {
	return r.visible[playlistID]
}

func (r *PlaylistPublishRepository) Commits() int { return r.commits }
