package model

import (
	"strings"
	"time"
)

const (
	ReportTargetVideo   = "video"
	ReportTargetComment = "comment"

	ReportStatusPending   = "pending"
	ReportStatusProcessed = "processed"
	ReportStatusRejected  = "rejected"
)

// Report 举报记录。
type Report struct {
	ID         string    `json:"id"`
	ReporterID string    `json:"reporter_id"`
	TargetType string    `json:"target_type"`
	TargetID   string    `json:"target_id"`
	Reason     string    `json:"reason"`
	Status     string    `json:"status"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Validate 校验并规范化举报字段。
func (r *Report) Validate() error {
	r.ReporterID = strings.TrimSpace(r.ReporterID)
	r.TargetType = strings.TrimSpace(r.TargetType)
	r.TargetID = strings.TrimSpace(r.TargetID)
	r.Reason = strings.TrimSpace(r.Reason)
	r.Note = strings.TrimSpace(r.Note)
	if r.ReporterID == "" {
		return NewValidationError("reporter_id", "举报人不能为空")
	}
	if r.TargetType == "" {
		r.TargetType = ReportTargetVideo
	}
	if r.TargetType != ReportTargetVideo && r.TargetType != ReportTargetComment {
		return NewValidationError("target_type", "举报对象类型不合法")
	}
	if r.TargetID == "" {
		return NewValidationError("target_id", "举报对象不能为空")
	}
	if r.Status == "" {
		r.Status = ReportStatusPending
	}
	if r.Status != ReportStatusPending && r.Status != ReportStatusProcessed && r.Status != ReportStatusRejected {
		return NewValidationError("status", "举报状态不合法")
	}
	return nil
}

var reportTransitions = map[string]map[string]bool{
	ReportStatusPending:   {ReportStatusProcessed: true, ReportStatusRejected: true},
	ReportStatusProcessed: {},
	ReportStatusRejected:  {},
}

// CanReportTransition 判断举报状态是否可流转。
func CanReportTransition(from, to string) bool {
	if m, ok := reportTransitions[from]; ok {
		return m[to]
	}
	return false
}

// ReportFilter 举报列表筛选条件。
type ReportFilter struct {
	TargetType string
	Status     string
	ReporterID string
}

// Match 判断举报是否命中筛选条件。
func (f ReportFilter) Match(r *Report) bool {
	if f.TargetType != "" && r.TargetType != f.TargetType {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.ReporterID != "" && r.ReporterID != f.ReporterID {
		return false
	}
	return true
}
