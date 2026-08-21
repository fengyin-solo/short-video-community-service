package service

import (
	"testing"

	"shortvideo/internal/config"
	"shortvideo/internal/model"
	"shortvideo/internal/store"
	"shortvideo/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func seedUser(t *testing.T, s *Service, username string) string {
	t.Helper()
	u, err := s.CreateUser(model.User{Username: username})
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return u.ID
}

func seedPublishedVideo(t *testing.T, s *Service, authorID string) string {
	t.Helper()
	v, err := s.CreateVideo(model.Video{Title: "测试视频", AuthorID: authorID, Duration: 60})
	if err != nil {
		t.Fatalf("seed video: %v", err)
	}
	if _, err := s.SubmitVideoForReview(v.ID); err != nil {
		t.Fatalf("submit: %v", err)
	}
	v, err = s.ApproveVideo(v.ID)
	if err != nil {
		t.Fatalf("approve: %v", err)
	}
	return v.ID
}

func TestVideoLifecycle(t *testing.T) {
	s := newTestService()
	authorID := seedUser(t, s, "alice")
	v, err := s.CreateVideo(model.Video{Title: "视频", AuthorID: authorID})
	if err != nil {
		t.Fatalf("create video: %v", err)
	}
	if v.Status != model.VideoStatusDraft {
		t.Fatalf("expect draft, got %s", v.Status)
	}
	v, _ = s.SubmitVideoForReview(v.ID)
	if v.Status != model.VideoStatusReviewing {
		t.Fatalf("expect reviewing, got %s", v.Status)
	}
	v, _ = s.ApproveVideo(v.ID)
	if v.Status != model.VideoStatusPublished {
		t.Fatalf("expect published, got %s", v.Status)
	}
	v, _ = s.BanVideo(v.ID)
	if v.Status != model.VideoStatusBanned {
		t.Fatalf("expect banned, got %s", v.Status)
	}
	// 已下架不能再审核
	if _, err := s.ApproveVideo(v.ID); err == nil {
		t.Fatal("expect conflict approving banned video")
	}
}

func TestLikeSideEffect(t *testing.T) {
	s := newTestService()
	authorID := seedUser(t, s, "alice")
	userID := seedUser(t, s, "bob")
	videoID := seedPublishedVideo(t, s, authorID)
	if _, err := s.LikeVideo(userID, videoID); err != nil {
		t.Fatalf("like: %v", err)
	}
	v, _ := s.GetVideo(videoID)
	if v.LikeCount != 1 {
		t.Fatalf("expect like count 1, got %d", v.LikeCount)
	}
	// 重复点赞应冲突
	if _, err := s.LikeVideo(userID, videoID); err == nil {
		t.Fatal("expect conflict duplicate like")
	}
	if err := s.UnlikeVideo(userID, videoID); err != nil {
		t.Fatalf("unlike: %v", err)
	}
	v, _ = s.GetVideo(videoID)
	if v.LikeCount != 0 {
		t.Fatalf("expect like count 0, got %d", v.LikeCount)
	}
}

func TestFollowSideEffect(t *testing.T) {
	s := newTestService()
	aID := seedUser(t, s, "alice")
	bID := seedUser(t, s, "bob")
	if _, err := s.FollowUser(aID, bID); err != nil {
		t.Fatalf("follow: %v", err)
	}
	a, _ := s.GetUser(aID)
	b, _ := s.GetUser(bID)
	if a.FollowingCount != 1 || b.FollowerCount != 1 {
		t.Fatalf("expect counts 1/1, got %d/%d", a.FollowingCount, b.FollowerCount)
	}
	if err := s.UnfollowUser(aID, bID); err != nil {
		t.Fatalf("unfollow: %v", err)
	}
	a, _ = s.GetUser(aID)
	b, _ = s.GetUser(bID)
	if a.FollowingCount != 0 || b.FollowerCount != 0 {
		t.Fatalf("expect counts 0/0, got %d/%d", a.FollowingCount, b.FollowerCount)
	}
}

func TestCommentSideEffect(t *testing.T) {
	s := newTestService()
	authorID := seedUser(t, s, "alice")
	userID := seedUser(t, s, "bob")
	videoID := seedPublishedVideo(t, s, authorID)
	c, err := s.CreateComment(model.Comment{VideoID: videoID, UserID: userID, Content: "很好"})
	if err != nil {
		t.Fatalf("create comment: %v", err)
	}
	v, _ := s.GetVideo(videoID)
	if v.CommentCount != 1 {
		t.Fatalf("expect comment count 1, got %d", v.CommentCount)
	}
	if _, err := s.DeleteComment(c.ID); err != nil {
		t.Fatalf("delete comment: %v", err)
	}
	v, _ = s.GetVideo(videoID)
	if v.CommentCount != 0 {
		t.Fatalf("expect comment count 0, got %d", v.CommentCount)
	}
}

func TestReportProcessBansVideo(t *testing.T) {
	s := newTestService()
	authorID := seedUser(t, s, "alice")
	reporterID := seedUser(t, s, "bob")
	videoID := seedPublishedVideo(t, s, authorID)
	rp, err := s.CreateReport(model.Report{
		ReporterID: reporterID, TargetType: model.ReportTargetVideo, TargetID: videoID, Reason: "违规",
	})
	if err != nil {
		t.Fatalf("create report: %v", err)
	}
	rp, err = s.ProcessReport(rp.ID, "确认违规")
	if err != nil {
		t.Fatalf("process report: %v", err)
	}
	if rp.Status != model.ReportStatusProcessed {
		t.Fatalf("expect processed, got %s", rp.Status)
	}
	v, _ := s.GetVideo(videoID)
	if v.Status != model.VideoStatusBanned {
		t.Fatalf("expect video banned, got %s", v.Status)
	}
}
