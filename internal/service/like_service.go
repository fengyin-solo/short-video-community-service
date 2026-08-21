package service

import (
	"sort"
	"time"

	"shortvideo/internal/model"
	"shortvideo/pkg/idgen"
)

// LikeVideo 点赞视频（联动视频点赞数 +1）。
func (s *Service) LikeVideo(userID, videoID string) (*model.Like, error) {
	if _, err := s.store.GetUser(userID); err != nil {
		return nil, model.NewValidationError("user_id", "用户不存在")
	}
	if _, err := s.store.GetVideo(videoID); err != nil {
		return nil, model.NewValidationError("video_id", "视频不存在")
	}
	like := &model.Like{ID: idgen.Hex(), UserID: userID, VideoID: videoID, CreatedAt: time.Now()}
	if err := like.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateLike(like); err != nil {
		return nil, err
	}
	if v, err := s.store.GetVideo(videoID); err == nil {
		v.LikeCount++
		v.UpdatedAt = time.Now()
		_ = s.store.UpdateVideo(v)
	}
	return like, nil
}

// UnlikeVideo 取消点赞（联动视频点赞数 -1）。
func (s *Service) UnlikeVideo(userID, videoID string) error {
	like, err := s.store.GetLikeByPair(userID, videoID)
	if err != nil {
		return err
	}
	if err := s.store.DeleteLike(like.ID); err != nil {
		return err
	}
	if v, err := s.store.GetVideo(videoID); err == nil {
		if v.LikeCount > 0 {
			v.LikeCount--
			v.UpdatedAt = time.Now()
			_ = s.store.UpdateVideo(v)
		}
	}
	return nil
}

// ListLikes 分页列出点赞记录。
func (s *Service) ListLikes(filter model.LikeFilter, page, size int) ([]*model.Like, int, error) {
	all := s.store.ListLikes()
	matched := make([]*model.Like, 0, len(all))
	for _, l := range all {
		if filter.Match(l) {
			matched = append(matched, l)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Like{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
