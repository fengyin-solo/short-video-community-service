package service

import (
	"sort"
	"time"

	"shortvideo/internal/model"
	"shortvideo/pkg/idgen"
)

// AttachTag 给视频打标签（联动标签关联视频数 +1）。
func (s *Service) AttachTag(videoID, tagID string) (*model.VideoTag, error) {
	if _, err := s.store.GetVideo(videoID); err != nil {
		return nil, model.NewValidationError("video_id", "视频不存在")
	}
	if _, err := s.store.GetTag(tagID); err != nil {
		return nil, model.NewValidationError("tag_id", "标签不存在")
	}
	vt := &model.VideoTag{ID: idgen.Hex(), VideoID: videoID, TagID: tagID, CreatedAt: time.Now()}
	if err := vt.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateVideoTag(vt); err != nil {
		return nil, err
	}
	if t, err := s.store.GetTag(tagID); err == nil {
		t.VideoCount++
		t.UpdatedAt = time.Now()
		_ = s.store.UpdateTag(t)
	}
	return vt, nil
}

// DetachTag 取消视频标签（联动标签关联视频数 -1）。
func (s *Service) DetachTag(videoID, tagID string) error {
	vt, err := s.store.GetVideoTagByPair(videoID, tagID)
	if err != nil {
		return err
	}
	if err := s.store.DeleteVideoTag(vt.ID); err != nil {
		return err
	}
	if t, err := s.store.GetTag(tagID); err == nil {
		if t.VideoCount > 0 {
			t.VideoCount--
			t.UpdatedAt = time.Now()
			_ = s.store.UpdateTag(t)
		}
	}
	return nil
}

// ListVideoTags 分页列出视频标签关联。
func (s *Service) ListVideoTags(filter model.VideoTagFilter, page, size int) ([]*model.VideoTag, int, error) {
	all := s.store.ListVideoTags()
	matched := make([]*model.VideoTag, 0, len(all))
	for _, v := range all {
		if filter.Match(v) {
			matched = append(matched, v)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.VideoTag{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// ListVideoTagNames 返回某视频关联的所有标签名。
func (s *Service) ListVideoTagNames(videoID string) []string {
	names := make([]string, 0)
	for _, vt := range s.store.ListVideoTags() {
		if vt.VideoID != videoID {
			continue
		}
		if t, err := s.store.GetTag(vt.TagID); err == nil {
			names = append(names, t.Name)
		}
	}
	return names
}
