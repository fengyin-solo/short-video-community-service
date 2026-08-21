// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"shortvideo/internal/model"
)

var (
	// ErrNotFound 表示记录不存在。
	ErrNotFound = errors.New("记录不存在")
	// ErrConflict 表示记录已存在或状态冲突。
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法。
type Store interface {
	CreateUser(u *model.User) error
	GetUser(id string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	ListUsers() []*model.User
	UpdateUser(u *model.User) error
	DeleteUser(id string) error

	CreateCategory(c *model.Category) error
	GetCategory(id string) (*model.Category, error)
	GetCategoryByName(name string) (*model.Category, error)
	ListCategories() []*model.Category
	UpdateCategory(c *model.Category) error
	DeleteCategory(id string) error

	CreateVideo(v *model.Video) error
	GetVideo(id string) (*model.Video, error)
	ListVideos() []*model.Video
	UpdateVideo(v *model.Video) error
	DeleteVideo(id string) error

	CreateTag(t *model.Tag) error
	GetTag(id string) (*model.Tag, error)
	GetTagByName(name string) (*model.Tag, error)
	ListTags() []*model.Tag
	UpdateTag(t *model.Tag) error
	DeleteTag(id string) error

	CreateLike(l *model.Like) error
	GetLike(id string) (*model.Like, error)
	GetLikeByPair(userID, videoID string) (*model.Like, error)
	ListLikes() []*model.Like
	DeleteLike(id string) error

	CreateComment(c *model.Comment) error
	GetComment(id string) (*model.Comment, error)
	ListComments() []*model.Comment
	UpdateComment(c *model.Comment) error
	DeleteComment(id string) error

	CreateFollow(f *model.Follow) error
	GetFollow(id string) (*model.Follow, error)
	GetFollowByPair(followerID, followeeID string) (*model.Follow, error)
	ListFollows() []*model.Follow
	DeleteFollow(id string) error

	CreateFavorite(f *model.Favorite) error
	GetFavorite(id string) (*model.Favorite, error)
	GetFavoriteByPair(userID, videoID string) (*model.Favorite, error)
	ListFavorites() []*model.Favorite
	DeleteFavorite(id string) error

	CreateReport(r *model.Report) error
	GetReport(id string) (*model.Report, error)
	ListReports() []*model.Report
	UpdateReport(r *model.Report) error
	DeleteReport(id string) error

	CreateNotification(n *model.Notification) error
	GetNotification(id string) (*model.Notification, error)
	ListNotifications() []*model.Notification
	UpdateNotification(n *model.Notification) error
	DeleteNotification(id string) error

	CreatePlaylist(p *model.Playlist) error
	GetPlaylist(id string) (*model.Playlist, error)
	ListPlaylists() []*model.Playlist
	UpdatePlaylist(p *model.Playlist) error
	DeletePlaylist(id string) error

	CreateVideoTag(v *model.VideoTag) error
	GetVideoTag(id string) (*model.VideoTag, error)
	GetVideoTagByPair(videoID, tagID string) (*model.VideoTag, error)
	ListVideoTags() []*model.VideoTag
	DeleteVideoTag(id string) error
}
