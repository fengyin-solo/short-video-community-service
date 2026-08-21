package service

import (
	"sort"
	"time"

	"shortvideo/internal/model"
	"shortvideo/pkg/idgen"
)

// AddFavorite 收藏视频。
func (s *Service) AddFavorite(userID, videoID, folder string) (*model.Favorite, error) {
	if _, err := s.store.GetUser(userID); err != nil {
		return nil, model.NewValidationError("user_id", "用户不存在")
	}
	if _, err := s.store.GetVideo(videoID); err != nil {
		return nil, model.NewValidationError("video_id", "视频不存在")
	}
	fav := &model.Favorite{ID: idgen.Hex(), UserID: userID, VideoID: videoID, Folder: folder, CreatedAt: time.Now()}
	if err := fav.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateFavorite(fav); err != nil {
		return nil, err
	}
	return fav, nil
}

// RemoveFavorite 取消收藏。
func (s *Service) RemoveFavorite(userID, videoID string) error {
	fav, err := s.store.GetFavoriteByPair(userID, videoID)
	if err != nil {
		return err
	}
	return s.store.DeleteFavorite(fav.ID)
}

// ListFavorites 分页列出收藏记录。
func (s *Service) ListFavorites(filter model.FavoriteFilter, page, size int) ([]*model.Favorite, int, error) {
	all := s.store.ListFavorites()
	matched := make([]*model.Favorite, 0, len(all))
	for _, f := range all {
		if filter.Match(f) {
			matched = append(matched, f)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Favorite{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
