package model

import (
	"strings"
	"time"
)

// VideoTag 视频-标签关联。
type VideoTag struct {
	ID        string    `json:"id"`
	VideoID   string    `json:"video_id"`
	TagID     string    `json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 校验视频标签关联。
func (v *VideoTag) Validate() error {
	v.VideoID = strings.TrimSpace(v.VideoID)
	v.TagID = strings.TrimSpace(v.TagID)
	if v.VideoID == "" {
		return NewValidationError("video_id", "视频不能为空")
	}
	if v.TagID == "" {
		return NewValidationError("tag_id", "标签不能为空")
	}
	return nil
}

// VideoTagFilter 视频标签列表筛选条件。
type VideoTagFilter struct {
	VideoID string
	TagID   string
}

// Match 判断关联是否命中筛选条件。
func (f VideoTagFilter) Match(v *VideoTag) bool {
	if f.VideoID != "" && v.VideoID != f.VideoID {
		return false
	}
	if f.TagID != "" && v.TagID != f.TagID {
		return false
	}
	return true
}
