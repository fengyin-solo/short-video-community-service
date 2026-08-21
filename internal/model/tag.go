package model

import (
	"strings"
	"time"
)

// Tag 视频标签。
type Tag struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	VideoCount int       `json:"video_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Validate 校验并规范化标签字段。
func (t *Tag) Validate() error {
	t.Name = strings.TrimSpace(t.Name)
	if t.Name == "" {
		return NewValidationError("name", "标签名称不能为空")
	}
	if t.VideoCount < 0 {
		return NewValidationError("video_count", "关联视频数不能为负")
	}
	return nil
}
