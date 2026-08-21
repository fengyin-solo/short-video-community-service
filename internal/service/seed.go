package service

import (
	"fmt"

	"shortvideo/internal/model"
)

// SeedResult 初始化演示数据的结果统计。
type SeedResult struct {
	Users      int `json:"users"`
	Categories int `json:"categories"`
	Videos     int `json:"videos"`
	Tags       int `json:"tags"`
	Likes      int `json:"likes"`
	Comments   int `json:"comments"`
	Follows    int `json:"follows"`
	Favorites  int `json:"favorites"`
	Reports    int `json:"reports"`
}

// SeedDemoData 初始化一批演示数据。
func (s *Service) SeedDemoData() (*SeedResult, error) {
	res := &SeedResult{}

	categories := []model.Category{
		{Name: "科技", Sort: 1}, {Name: "美食", Sort: 2}, {Name: "旅行", Sort: 3},
		{Name: "运动", Sort: 4}, {Name: "音乐", Sort: 5},
	}
	catByName := make(map[string]string)
	for _, c := range categories {
		created, err := s.CreateCategory(c)
		if err != nil {
			return nil, fmt.Errorf("seed category: %w", err)
		}
		catByName[c.Name] = created.ID
		res.Categories++
	}

	users := []model.User{
		{Username: "alice", Nickname: "爱丽丝", Bio: "科技博主"},
		{Username: "bob", Nickname: "鲍勃", Bio: "美食探店"},
		{Username: "carol", Nickname: "卡罗尔", Bio: "旅行达人"},
		{Username: "dave", Nickname: "戴夫", Bio: "运动健身"},
		{Username: "erin", Nickname: "艾琳", Bio: "音乐人"},
	}
	userByName := make(map[string]string)
	for _, u := range users {
		created, err := s.CreateUser(u)
		if err != nil {
			return nil, fmt.Errorf("seed user: %w", err)
		}
		userByName[u.Username] = created.ID
		res.Users++
	}

	// 视频
	videos := []struct {
		title, author, category string
		views                   int
	}{
		{"新一代芯片性能评测", "alice", "科技", 52000},
		{"五分钟看懂量子计算", "alice", "科技", 31000},
		{"街头美食探店第一弹", "bob", "美食", 88000},
		{"在家做正宗红烧肉", "bob", "美食", 45000},
		{"川西环线自驾攻略", "carol", "旅行", 67000},
		{"海岛潜水 vlog", "carol", "旅行", 23000},
		{"核心力量训练计划", "dave", "运动", 15000},
		{"吉他弹唱翻唱", "erin", "音乐", 9000},
	}
	videoIDs := make([]string, 0, len(videos))
	for _, v := range videos {
		created, err := s.CreateVideo(model.Video{
			Title: v.title, AuthorID: userByName[v.author], CategoryID: catByName[v.category],
			CoverURL: "https://example.com/cover.jpg", VideoURL: "https://example.com/v.mp4", Duration: 60,
		})
		if err != nil {
			return nil, fmt.Errorf("seed video: %w", err)
		}
		if _, err := s.SubmitVideoForReview(created.ID); err != nil {
			return nil, err
		}
		created, err = s.ApproveVideo(created.ID)
		if err != nil {
			return nil, err
		}
		created.ViewCount = v.views
		_ = s.store.UpdateVideo(created)
		videoIDs = append(videoIDs, created.ID)
		res.Videos++
	}

	// 标签
	tags := []string{"芯片", "美食", "旅行", "健身", "音乐"}
	for _, t := range tags {
		if _, err := s.CreateTag(model.Tag{Name: t}); err != nil {
			return nil, err
		}
		res.Tags++
	}

	// 点赞
	if _, err := s.LikeVideo(userByName["bob"], videoIDs[0]); err != nil {
		return nil, err
	}
	if _, err := s.LikeVideo(userByName["carol"], videoIDs[0]); err != nil {
		return nil, err
	}
	if _, err := s.LikeVideo(userByName["alice"], videoIDs[2]); err != nil {
		return nil, err
	}
	res.Likes = 3

	// 评论
	if _, err := s.CreateComment(model.Comment{VideoID: videoIDs[0], UserID: userByName["bob"], Content: "讲得太清楚了！"}); err != nil {
		return nil, err
	}
	if _, err := s.CreateComment(model.Comment{VideoID: videoIDs[2], UserID: userByName["alice"], Content: "看着就饿了"}); err != nil {
		return nil, err
	}
	res.Comments = 2

	// 关注
	if _, err := s.FollowUser(userByName["bob"], userByName["alice"]); err != nil {
		return nil, err
	}
	if _, err := s.FollowUser(userByName["carol"], userByName["alice"]); err != nil {
		return nil, err
	}
	if _, err := s.FollowUser(userByName["alice"], userByName["carol"]); err != nil {
		return nil, err
	}
	res.Follows = 3

	// 收藏
	if _, err := s.AddFavorite(userByName["bob"], videoIDs[0], "默认收藏夹"); err != nil {
		return nil, err
	}
	res.Favorites = 1

	// 举报
	if _, err := s.CreateReport(model.Report{
		ReporterID: userByName["dave"], TargetType: model.ReportTargetVideo, TargetID: videoIDs[1], Reason: "内容涉嫌误导",
	}); err != nil {
		return nil, err
	}
	res.Reports = 1

	return res, nil
}
