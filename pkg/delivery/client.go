package delivery

import "errors"

var (
	ErrRejected  = errors.New("source rejected")
	ErrTemporary = errors.New("source temporarily unavailable")
)

type Client struct {
	Calls     map[string]int
	Published map[string]int
}

func NewClient() *Client { return &Client{Calls: map[string]int{}, Published: map[string]int{}} }

func (c *Client) Send(videoID string) error {
	c.Calls[videoID]++
	switch videoID {
	case "video-rejected":
		return ErrRejected
	case "video-temporary":
		if c.Calls[videoID] == 1 {
			return ErrTemporary
		}
	}
	c.Published[videoID]++
	return nil
}
