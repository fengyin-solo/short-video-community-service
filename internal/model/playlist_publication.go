package model

import "fmt"

type PlaylistPublication struct {
	PlaylistID string
	Attempt    int
	Status     string
	Failure    string
}

func (p PlaylistPublication) EventKey() string {
	return fmt.Sprintf("playlist:%s:published:%d", p.PlaylistID, p.Attempt)
}
