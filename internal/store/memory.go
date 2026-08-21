package store

import (
	"sync"

	"shortvideo/internal/model"
)

// MemoryStore 基于内存的 Store 实现。
type MemoryStore struct {
	mu sync.RWMutex

	users      map[string]*model.User
	categories map[string]*model.Category
	videos     map[string]*model.Video
	tags       map[string]*model.Tag
	likes      map[string]*model.Like
	comments   map[string]*model.Comment
	follows    map[string]*model.Follow
	favorites  map[string]*model.Favorite
	reports    map[string]*model.Report

	notifications map[string]*model.Notification
	playlists     map[string]*model.Playlist
	videoTags     map[string]*model.VideoTag
}

// NewMemoryStore 构造内存存储。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:      make(map[string]*model.User),
		categories: make(map[string]*model.Category),
		videos:     make(map[string]*model.Video),
		tags:       make(map[string]*model.Tag),
		likes:      make(map[string]*model.Like),
		comments:   make(map[string]*model.Comment),
		follows:    make(map[string]*model.Follow),
		favorites:  make(map[string]*model.Favorite),
		reports:    make(map[string]*model.Report),

		notifications: make(map[string]*model.Notification),
		playlists:     make(map[string]*model.Playlist),
		videoTags:     make(map[string]*model.VideoTag),
	}
}

var _ Store = (*MemoryStore)(nil)
