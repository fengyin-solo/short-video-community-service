package service

import (
	"sort"
	"time"

	"shortvideo/internal/model"
	"shortvideo/pkg/idgen"
)

// FollowUser 关注用户（联动双方计数）。
func (s *Service) FollowUser(followerID, followeeID string) (*model.Follow, error) {
	follow := &model.Follow{ID: idgen.Hex(), FollowerID: followerID, FolloweeID: followeeID, CreatedAt: time.Now()}
	if err := follow.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetUser(followerID); err != nil {
		return nil, model.NewValidationError("follower_id", "关注者不存在")
	}
	if _, err := s.store.GetUser(followeeID); err != nil {
		return nil, model.NewValidationError("followee_id", "被关注者不存在")
	}
	if err := s.store.CreateFollow(follow); err != nil {
		return nil, err
	}
	if f, err := s.store.GetUser(followerID); err == nil {
		f.FollowingCount++
		f.UpdatedAt = time.Now()
		_ = s.store.UpdateUser(f)
	}
	if f, err := s.store.GetUser(followeeID); err == nil {
		f.FollowerCount++
		f.UpdatedAt = time.Now()
		_ = s.store.UpdateUser(f)
	}
	return follow, nil
}

// UnfollowUser 取消关注（联动双方计数）。
func (s *Service) UnfollowUser(followerID, followeeID string) error {
	follow, err := s.store.GetFollowByPair(followerID, followeeID)
	if err != nil {
		return err
	}
	if err := s.store.DeleteFollow(follow.ID); err != nil {
		return err
	}
	if f, err := s.store.GetUser(followerID); err == nil {
		if f.FollowingCount > 0 {
			f.FollowingCount--
			f.UpdatedAt = time.Now()
			_ = s.store.UpdateUser(f)
		}
	}
	if f, err := s.store.GetUser(followeeID); err == nil {
		if f.FollowerCount > 0 {
			f.FollowerCount--
			f.UpdatedAt = time.Now()
			_ = s.store.UpdateUser(f)
		}
	}
	return nil
}

// ListFollows 分页列出关注关系。
func (s *Service) ListFollows(filter model.FollowFilter, page, size int) ([]*model.Follow, int, error) {
	all := s.store.ListFollows()
	matched := make([]*model.Follow, 0, len(all))
	for _, f := range all {
		if filter.Match(f) {
			matched = append(matched, f)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Follow{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
