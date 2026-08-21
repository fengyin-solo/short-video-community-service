package store

import (
	"testing"
	"time"

	"shortvideo/internal/model"
)

func TestUserCRUD(t *testing.T) {
	s := NewMemoryStore()
	u := &model.User{ID: "u1", Username: "alice", Status: model.UserActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateUser(u); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateUser(&model.User{ID: "u2", Username: "alice"}); err != ErrConflict {
		t.Fatalf("expect conflict, got %v", err)
	}
	byName, err := s.GetUserByUsername("alice")
	if err != nil || byName.ID != "u1" {
		t.Fatalf("get by username failed: %v", err)
	}
	u.Nickname = "爱丽丝"
	if err := s.UpdateUser(u); err != nil {
		t.Fatalf("update: %v", err)
	}
	if len(s.ListUsers()) != 1 {
		t.Fatalf("expect 1 user")
	}
	if err := s.DeleteUser("u1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestCategoryCRUD(t *testing.T) {
	s := NewMemoryStore()
	c := &model.Category{ID: "c1", Name: "科技", Status: model.CategoryActive}
	if err := s.CreateCategory(c); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateCategory(&model.Category{ID: "c2", Name: "科技"}); err != ErrConflict {
		t.Fatalf("expect conflict, got %v", err)
	}
	if _, err := s.GetCategoryByName("科技"); err != nil {
		t.Fatalf("get by name: %v", err)
	}
	if len(s.ListCategories()) != 1 {
		t.Fatalf("expect 1 category")
	}
	if err := s.DeleteCategory("c1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestVideoCRUD(t *testing.T) {
	s := NewMemoryStore()
	v := &model.Video{ID: "v1", Title: "测试视频", AuthorID: "u1", Status: model.VideoStatusDraft}
	if err := s.CreateVideo(v); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := s.GetVideo("v1")
	if err != nil || got.Title != "测试视频" {
		t.Fatalf("get failed: %v", err)
	}
	got.Status = model.VideoStatusPublished
	if err := s.UpdateVideo(got); err != nil {
		t.Fatalf("update: %v", err)
	}
	if len(s.ListVideos()) != 1 {
		t.Fatalf("expect 1 video")
	}
	if err := s.DeleteVideo("v1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestLikeUniqueness(t *testing.T) {
	s := NewMemoryStore()
	l := &model.Like{ID: "l1", UserID: "u1", VideoID: "v1"}
	if err := s.CreateLike(l); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateLike(&model.Like{ID: "l2", UserID: "u1", VideoID: "v1"}); err != ErrConflict {
		t.Fatalf("expect conflict, got %v", err)
	}
	byPair, err := s.GetLikeByPair("u1", "v1")
	if err != nil || byPair.ID != "l1" {
		t.Fatalf("get by pair failed: %v", err)
	}
	if err := s.DeleteLike("l1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestCommentCRUD(t *testing.T) {
	s := NewMemoryStore()
	c := &model.Comment{ID: "c1", VideoID: "v1", UserID: "u1", Content: "不错", Status: model.CommentNormal}
	if err := s.CreateComment(c); err != nil {
		t.Fatalf("create: %v", err)
	}
	c.Status = model.CommentDeleted
	if err := s.UpdateComment(c); err != nil {
		t.Fatalf("update: %v", err)
	}
	if len(s.ListComments()) != 1 {
		t.Fatalf("expect 1 comment")
	}
	if err := s.DeleteComment("c1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestFollowUniqueness(t *testing.T) {
	s := NewMemoryStore()
	f := &model.Follow{ID: "f1", FollowerID: "u1", FolloweeID: "u2"}
	if err := s.CreateFollow(f); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateFollow(&model.Follow{ID: "f2", FollowerID: "u1", FolloweeID: "u2"}); err != ErrConflict {
		t.Fatalf("expect conflict, got %v", err)
	}
	if _, err := s.GetFollowByPair("u1", "u2"); err != nil {
		t.Fatalf("get by pair: %v", err)
	}
	if err := s.DeleteFollow("f1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestFavoriteUniqueness(t *testing.T) {
	s := NewMemoryStore()
	fav := &model.Favorite{ID: "f1", UserID: "u1", VideoID: "v1"}
	if err := s.CreateFavorite(fav); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateFavorite(&model.Favorite{ID: "f2", UserID: "u1", VideoID: "v1"}); err != ErrConflict {
		t.Fatalf("expect conflict, got %v", err)
	}
	if _, err := s.GetFavoriteByPair("u1", "v1"); err != nil {
		t.Fatalf("get by pair: %v", err)
	}
	if err := s.DeleteFavorite("f1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestReportCRUD(t *testing.T) {
	s := NewMemoryStore()
	r := &model.Report{ID: "r1", ReporterID: "u1", TargetType: model.ReportTargetVideo, TargetID: "v1", Status: model.ReportStatusPending}
	if err := s.CreateReport(r); err != nil {
		t.Fatalf("create: %v", err)
	}
	r.Status = model.ReportStatusProcessed
	if err := s.UpdateReport(r); err != nil {
		t.Fatalf("update: %v", err)
	}
	if len(s.ListReports()) != 1 {
		t.Fatalf("expect 1 report")
	}
	if err := s.DeleteReport("r1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestTagCRUD(t *testing.T) {
	s := NewMemoryStore()
	tg := &model.Tag{ID: "t1", Name: "芯片"}
	if err := s.CreateTag(tg); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateTag(&model.Tag{ID: "t2", Name: "芯片"}); err != ErrConflict {
		t.Fatalf("expect conflict, got %v", err)
	}
	if _, err := s.GetTagByName("芯片"); err != nil {
		t.Fatalf("get by name: %v", err)
	}
	if len(s.ListTags()) != 1 {
		t.Fatalf("expect 1 tag")
	}
	if err := s.DeleteTag("t1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
