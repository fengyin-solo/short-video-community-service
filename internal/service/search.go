package service

import (
	"sort"
	"strings"

	"shortvideo/internal/model"
)

// SearchVideos 按关键词搜索已发布视频（匹配标题与描述），按相关度排序。
func (s *Service) SearchVideos(query, categoryID string, page, size int) ([]*model.Video, int, error) {
	query = strings.ToLower(strings.TrimSpace(query))
	all := s.store.ListVideos()
	type scored struct {
		v     *model.Video
		score int
	}
	matched := make([]scored, 0)
	for _, v := range all {
		if v.Status != model.VideoStatusPublished {
			continue
		}
		if categoryID != "" && v.CategoryID != categoryID {
			continue
		}
		if query == "" {
			matched = append(matched, scored{v: v, score: v.ViewCount})
			continue
		}
		title := strings.ToLower(v.Title)
		desc := strings.ToLower(v.Description)
		score := 0
		if title == query {
			score += 100
		} else if strings.Contains(title, query) {
			score += 50
		}
		if strings.Contains(desc, query) {
			score += 20
		}
		if score > 0 {
			matched = append(matched, scored{v: v, score: score})
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].score != matched[j].score {
			return matched[i].score > matched[j].score
		}
		return matched[i].v.ViewCount > matched[j].v.ViewCount
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Video{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	result := make([]*model.Video, 0, end-start)
	for _, sc := range matched[start:end] {
		result = append(result, sc.v)
	}
	return result, total, nil
}

// FeedVideos 返回已发布视频的信息流（按播放量倒序）。
func (s *Service) FeedVideos(page, size int) ([]*model.Video, int, error) {
	all := s.store.ListVideos()
	published := make([]*model.Video, 0, len(all))
	for _, v := range all {
		if v.Status == model.VideoStatusPublished {
			published = append(published, v)
		}
	}
	sort.Slice(published, func(i, j int) bool { return published[i].ViewCount > published[j].ViewCount })
	total := len(published)
	start := (page - 1) * size
	if start >= total {
		return []*model.Video{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return published[start:end], total, nil
}

// ListVideosByAuthor 返回某作者的全部视频。
func (s *Service) ListVideosByAuthor(authorID string, page, size int) ([]*model.Video, int, error) {
	all := s.store.ListVideos()
	matched := make([]*model.Video, 0, len(all))
	for _, v := range all {
		if v.AuthorID == authorID {
			matched = append(matched, v)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Video{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
