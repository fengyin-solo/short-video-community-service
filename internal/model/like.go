package model

import (
	"strings"
	"time"
)

// Like 点赞记录。
type Like struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	VideoID   string    `json:"video_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 校验点赞记录。
func (l *Like) Validate() error {
	l.UserID = strings.TrimSpace(l.UserID)
	l.VideoID = strings.TrimSpace(l.VideoID)
	if l.UserID == "" {
		return NewValidationError("user_id", "用户不能为空")
	}
	if l.VideoID == "" {
		return NewValidationError("video_id", "视频不能为空")
	}
	return nil
}

// LikeFilter 点赞列表筛选条件。
type LikeFilter struct {
	UserID  string
	VideoID string
}

// Match 判断点赞记录是否命中筛选条件。
func (f LikeFilter) Match(l *Like) bool {
	if f.UserID != "" && l.UserID != f.UserID {
		return false
	}
	if f.VideoID != "" && l.VideoID != f.VideoID {
		return false
	}
	return true
}
