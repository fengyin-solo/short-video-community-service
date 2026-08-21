package model

import (
	"strings"
	"time"
)

// Favorite 收藏记录。
type Favorite struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	VideoID   string    `json:"video_id"`
	Folder    string    `json:"folder"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 校验收藏记录。
func (f *Favorite) Validate() error {
	f.UserID = strings.TrimSpace(f.UserID)
	f.VideoID = strings.TrimSpace(f.VideoID)
	f.Folder = strings.TrimSpace(f.Folder)
	if f.UserID == "" {
		return NewValidationError("user_id", "用户不能为空")
	}
	if f.VideoID == "" {
		return NewValidationError("video_id", "视频不能为空")
	}
	return nil
}

// FavoriteFilter 收藏列表筛选条件。
type FavoriteFilter struct {
	UserID  string
	VideoID string
	Folder  string
}

// Match 判断收藏记录是否命中筛选条件。
func (f FavoriteFilter) Match(v *Favorite) bool {
	if f.UserID != "" && v.UserID != f.UserID {
		return false
	}
	if f.VideoID != "" && v.VideoID != f.VideoID {
		return false
	}
	if f.Folder != "" && v.Folder != f.Folder {
		return false
	}
	return true
}
