package playlistpublish

import "shortvideo/internal/model"

type Notification struct {
	PlaylistID string
	Attempt    int
	Premature  bool
}

type Publisher struct {
	notifications []Notification
	published     map[string]struct{}
}

func NewPublisher() *Publisher { return &Publisher{published: make(map[string]struct{})} }

func (p *Publisher) Publish(publication model.PlaylistPublication, visible bool) {
	key := publication.EventKey()
	if _, exists := p.published[key]; exists {
		return
	}
	p.published[key] = struct{}{}
	p.notifications = append(p.notifications, Notification{
		PlaylistID: publication.PlaylistID,
		Attempt:    publication.Attempt,
		Premature:  !visible,
	})
}

func (p *Publisher) Notifications() []Notification {
	return append([]Notification(nil), p.notifications...)
}
