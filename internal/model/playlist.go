package model

import (
	"strings"
	"time"
)

const (
	PlaylistPublic  = "public"
	PlaylistPrivate = "private"
)

// Playlist 用户播单。
type Playlist struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	VideoIDs    []string  `json:"video_ids"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate 校验并规范化播单字段。
func (p *Playlist) Validate() error {
	p.UserID = strings.TrimSpace(p.UserID)
	p.Name = strings.TrimSpace(p.Name)
	p.Description = strings.TrimSpace(p.Description)
	if p.UserID == "" {
		return NewValidationError("user_id", "用户不能为空")
	}
	if p.Name == "" {
		return NewValidationError("name", "播单名称不能为空")
	}
	if p.Status == "" {
		p.Status = PlaylistPublic
	}
	if p.Status != PlaylistPublic && p.Status != PlaylistPrivate {
		return NewValidationError("status", "播单状态不合法")
	}
	// 去重且保持顺序
	seen := make(map[string]bool)
	ids := make([]string, 0, len(p.VideoIDs))
	for _, id := range p.VideoIDs {
		id = strings.TrimSpace(id)
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		ids = append(ids, id)
	}
	p.VideoIDs = ids
	return nil
}

// PlaylistFilter 播单列表筛选条件。
type PlaylistFilter struct {
	UserID string
	Status string
}

// Match 判断播单是否命中筛选条件。
func (f PlaylistFilter) Match(p *Playlist) bool {
	if f.UserID != "" && p.UserID != f.UserID {
		return false
	}
	if f.Status != "" && p.Status != f.Status {
		return false
	}
	return true
}
