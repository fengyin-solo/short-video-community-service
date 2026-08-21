package model

import (
	"strings"
	"time"
)

const (
	CommentNormal  = "normal"
	CommentDeleted = "deleted"
)

// Comment 视频评论。
type Comment struct {
	ID        string    `json:"id"`
	VideoID   string    `json:"video_id"`
	UserID    string    `json:"user_id"`
	ParentID  string    `json:"parent_id"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 校验并规范化评论字段。
func (c *Comment) Validate() error {
	c.VideoID = strings.TrimSpace(c.VideoID)
	c.UserID = strings.TrimSpace(c.UserID)
	c.ParentID = strings.TrimSpace(c.ParentID)
	c.Content = strings.TrimSpace(c.Content)
	if c.VideoID == "" {
		return NewValidationError("video_id", "视频不能为空")
	}
	if c.UserID == "" {
		return NewValidationError("user_id", "用户不能为空")
	}
	if c.Content == "" {
		return NewValidationError("content", "评论内容不能为空")
	}
	if c.Status == "" {
		c.Status = CommentNormal
	}
	if c.Status != CommentNormal && c.Status != CommentDeleted {
		return NewValidationError("status", "评论状态不合法")
	}
	return nil
}

// CommentFilter 评论列表筛选条件。
type CommentFilter struct {
	VideoID string
	UserID  string
	Status  string
}

// Match 判断评论是否命中筛选条件。
func (f CommentFilter) Match(c *Comment) bool {
	if f.VideoID != "" && c.VideoID != f.VideoID {
		return false
	}
	if f.UserID != "" && c.UserID != f.UserID {
		return false
	}
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	return true
}
