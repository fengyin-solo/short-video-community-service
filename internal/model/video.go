package model

import (
	"strings"
	"time"
)

const (
	VideoStatusDraft     = "draft"
	VideoStatusReviewing = "reviewing"
	VideoStatusPublished = "published"
	VideoStatusBanned    = "banned"
)

// Video 视频。
type Video struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CoverURL     string    `json:"cover_url"`
	VideoURL     string    `json:"video_url"`
	AuthorID     string    `json:"author_id"`
	CategoryID   string    `json:"category_id"`
	Duration     int       `json:"duration"`
	Status       string    `json:"status"`
	ViewCount    int       `json:"view_count"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Validate 校验并规范化视频字段。
func (v *Video) Validate() error {
	v.Title = strings.TrimSpace(v.Title)
	v.Description = strings.TrimSpace(v.Description)
	v.CoverURL = strings.TrimSpace(v.CoverURL)
	v.VideoURL = strings.TrimSpace(v.VideoURL)
	v.AuthorID = strings.TrimSpace(v.AuthorID)
	v.CategoryID = strings.TrimSpace(v.CategoryID)
	if v.Title == "" {
		return NewValidationError("title", "视频标题不能为空")
	}
	if v.AuthorID == "" {
		return NewValidationError("author_id", "作者不能为空")
	}
	if v.Status == "" {
		v.Status = VideoStatusDraft
	}
	if !isVideoStatus(v.Status) {
		return NewValidationError("status", "视频状态不合法")
	}
	if v.Duration < 0 {
		return NewValidationError("duration", "视频时长不能为负")
	}
	if v.ViewCount < 0 || v.LikeCount < 0 || v.CommentCount < 0 {
		return NewValidationError("count", "计数不能为负")
	}
	return nil
}

func isVideoStatus(s string) bool {
	switch s {
	case VideoStatusDraft, VideoStatusReviewing, VideoStatusPublished, VideoStatusBanned:
		return true
	}
	return false
}

var videoTransitions = map[string]map[string]bool{
	VideoStatusDraft:     {VideoStatusReviewing: true},
	VideoStatusReviewing: {VideoStatusPublished: true, VideoStatusBanned: true},
	VideoStatusPublished: {VideoStatusBanned: true},
	VideoStatusBanned:    {},
}

// CanVideoTransition 判断视频状态是否可流转。
func CanVideoTransition(from, to string) bool {
	if m, ok := videoTransitions[from]; ok {
		return m[to]
	}
	return false
}

// VideoFilter 视频列表筛选条件。
type VideoFilter struct {
	Status     string
	CategoryID string
	AuthorID   string
	Keyword    string
}

// Match 判断视频是否命中筛选条件。
func (f VideoFilter) Match(v *Video) bool {
	if f.Status != "" && v.Status != f.Status {
		return false
	}
	if f.CategoryID != "" && v.CategoryID != f.CategoryID {
		return false
	}
	if f.AuthorID != "" && v.AuthorID != f.AuthorID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(v.Title), k) &&
			!strings.Contains(strings.ToLower(v.Description), k) {
			return false
		}
	}
	return true
}
