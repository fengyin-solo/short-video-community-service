package service

import (
	"sort"
	"time"

	"shortvideo/internal/model"
	"shortvideo/internal/store"
	"shortvideo/pkg/idgen"
)

// CreateReport 创建举报。
func (s *Service) CreateReport(input model.Report) (*model.Report, error) {
	input.ID = idgen.Hex()
	input.Status = model.ReportStatusPending
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if _, err := s.store.GetUser(input.ReporterID); err != nil {
		return nil, model.NewValidationError("reporter_id", "举报人不存在")
	}
	if input.TargetType == model.ReportTargetVideo {
		if _, err := s.store.GetVideo(input.TargetID); err != nil {
			return nil, model.NewValidationError("target_id", "被举报视频不存在")
		}
	} else {
		if _, err := s.store.GetComment(input.TargetID); err != nil {
			return nil, model.NewValidationError("target_id", "被举报评论不存在")
		}
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateReport(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

// GetReport 获取举报。
func (s *Service) GetReport(id string) (*model.Report, error) {
	return s.store.GetReport(id)
}

// ListReports 分页列出举报。
func (s *Service) ListReports(filter model.ReportFilter, page, size int) ([]*model.Report, int, error) {
	all := s.store.ListReports()
	matched := make([]*model.Report, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Report{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// ProcessReport 处理举报（有效举报，可联动下架目标视频）。
func (s *Service) ProcessReport(id, note string) (*model.Report, error) {
	existing, err := s.store.GetReport(id)
	if err != nil {
		return nil, err
	}
	if !model.CanReportTransition(existing.Status, model.ReportStatusProcessed) {
		return nil, store.ErrConflict
	}
	existing.Status = model.ReportStatusProcessed
	existing.Note = note
	existing.UpdatedAt = time.Now()
	if existing.TargetType == model.ReportTargetVideo {
		if v, err := s.store.GetVideo(existing.TargetID); err == nil {
			v.Status = model.VideoStatusBanned
			v.UpdatedAt = time.Now()
			_ = s.store.UpdateVideo(v)
		}
	}
	if err := s.store.UpdateReport(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// RejectReport 驳回举报。
func (s *Service) RejectReport(id, note string) (*model.Report, error) {
	existing, err := s.store.GetReport(id)
	if err != nil {
		return nil, err
	}
	if !model.CanReportTransition(existing.Status, model.ReportStatusRejected) {
		return nil, store.ErrConflict
	}
	existing.Status = model.ReportStatusRejected
	existing.Note = note
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateReport(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteReport 删除举报。
func (s *Service) DeleteReport(id string) error {
	return s.store.DeleteReport(id)
}
