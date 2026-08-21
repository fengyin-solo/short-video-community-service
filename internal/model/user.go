package model

import (
	"strings"
	"time"
)

const (
	UserActive = "active"
	UserBanned = "banned"
)

// User 用户。
type User struct {
	ID             string    `json:"id"`
	Username       string    `json:"username"`
	Nickname       string    `json:"nickname"`
	Avatar         string    `json:"avatar"`
	Bio            string    `json:"bio"`
	Status         string    `json:"status"`
	FollowerCount  int       `json:"follower_count"`
	FollowingCount int       `json:"following_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Validate 校验并规范化用户字段。
func (u *User) Validate() error {
	u.Username = strings.TrimSpace(u.Username)
	u.Nickname = strings.TrimSpace(u.Nickname)
	u.Avatar = strings.TrimSpace(u.Avatar)
	u.Bio = strings.TrimSpace(u.Bio)
	if u.Username == "" {
		return NewValidationError("username", "用户名不能为空")
	}
	if u.Nickname == "" {
		u.Nickname = u.Username
	}
	if u.Status == "" {
		u.Status = UserActive
	}
	if u.Status != UserActive && u.Status != UserBanned {
		return NewValidationError("status", "用户状态不合法")
	}
	if u.FollowerCount < 0 || u.FollowingCount < 0 {
		return NewValidationError("count", "关注/粉丝数不能为负")
	}
	return nil
}

// UserFilter 用户列表筛选条件。
type UserFilter struct {
	Status  string
	Keyword string
}

// Match 判断用户是否命中筛选条件。
func (f UserFilter) Match(u *User) bool {
	if f.Status != "" && u.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(u.Username), k) &&
			!strings.Contains(strings.ToLower(u.Nickname), k) {
			return false
		}
	}
	return true
}
