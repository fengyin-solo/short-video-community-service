package service

import (
	"sort"
	"time"

	"shortvideo/internal/model"
	"shortvideo/pkg/idgen"
)

// CreateNotification 创建通知。
func (s *Service) CreateNotification(input model.Notification) (*model.Notification, error) {
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateNotification(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

// notify 内部便捷方法：向指定用户推送一条通知。
func (s *Service) notify(userID, ntype, refID, content string) {
	n := &model.Notification{
		ID: idgen.Hex(), UserID: userID, Type: ntype, RefID: refID,
		Content: content, CreatedAt: time.Now(),
	}
	_ = s.store.CreateNotification(n)
}

// GetNotification 获取通知。
func (s *Service) GetNotification(id string) (*model.Notification, error) {
	return s.store.GetNotification(id)
}

// ListNotifications 分页列出通知。
func (s *Service) ListNotifications(filter model.NotificationFilter, page, size int) ([]*model.Notification, int, error) {
	all := s.store.ListNotifications()
	matched := make([]*model.Notification, 0, len(all))
	for _, n := range all {
		if filter.Match(n) {
			matched = append(matched, n)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Notification{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// MarkNotificationRead 标记单条通知为已读。
func (s *Service) MarkNotificationRead(id string) (*model.Notification, error) {
	n, err := s.store.GetNotification(id)
	if err != nil {
		return nil, err
	}
	n.Read = true
	if err := s.store.UpdateNotification(n); err != nil {
		return nil, err
	}
	return n, nil
}

// MarkAllNotificationsRead 将某用户全部通知标记为已读。
func (s *Service) MarkAllNotificationsRead(userID string) (int, error) {
	count := 0
	for _, n := range s.store.ListNotifications() {
		if n.UserID == userID && !n.Read {
			n.Read = true
			if err := s.store.UpdateNotification(n); err != nil {
				return count, err
			}
			count++
		}
	}
	return count, nil
}

// DeleteNotification 删除通知。
func (s *Service) DeleteNotification(id string) error {
	return s.store.DeleteNotification(id)
}
