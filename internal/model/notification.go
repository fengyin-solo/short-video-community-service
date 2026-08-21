package model

import (
	"strings"
	"time"
)

const (
	NotificationTypeLike    = "like"
	NotificationTypeComment = "comment"
	NotificationTypeFollow  = "follow"
	NotificationTypeSystem  = "system"
)

// Notification 站内通知。
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	RefID     string    `json:"ref_id"`
	Content   string    `json:"content"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 校验并规范化通知字段。
func (n *Notification) Validate() error {
	n.UserID = strings.TrimSpace(n.UserID)
	n.Type = strings.TrimSpace(n.Type)
	n.RefID = strings.TrimSpace(n.RefID)
	n.Content = strings.TrimSpace(n.Content)
	if n.UserID == "" {
		return NewValidationError("user_id", "接收人不能为空")
	}
	if n.Type == "" {
		n.Type = NotificationTypeSystem
	}
	switch n.Type {
	case NotificationTypeLike, NotificationTypeComment, NotificationTypeFollow, NotificationTypeSystem:
	default:
		return NewValidationError("type", "通知类型不合法")
	}
	if n.Content == "" {
		return NewValidationError("content", "通知内容不能为空")
	}
	return nil
}

// NotificationFilter 通知列表筛选条件。
type NotificationFilter struct {
	UserID string
	Type   string
	Read   *bool
}

// Match 判断通知是否命中筛选条件。
func (f NotificationFilter) Match(n *Notification) bool {
	if f.UserID != "" && n.UserID != f.UserID {
		return false
	}
	if f.Type != "" && n.Type != f.Type {
		return false
	}
	if f.Read != nil && n.Read != *f.Read {
		return false
	}
	return true
}
