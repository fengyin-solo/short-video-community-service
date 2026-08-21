package handler

import (
	"net/http"

	"shortvideo/internal/model"
	"shortvideo/pkg/httpx"
)

func (s *Server) registerReportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/reports", s.createReport)
	mux.HandleFunc("GET /api/reports", s.listReports)
	mux.HandleFunc("GET /api/reports/{id}", s.getReport)
	mux.HandleFunc("POST /api/reports/{id}/process", s.processReport)
	mux.HandleFunc("POST /api/reports/{id}/reject", s.rejectReport)
	mux.HandleFunc("DELETE /api/reports/{id}", s.deleteReport)
}

type reportRequest struct {
	ReporterID string `json:"reporter_id"`
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
	Reason     string `json:"reason"`
}

type reportActionRequest struct {
	Note string `json:"note"`
}

func (s *Server) createReport(w http.ResponseWriter, r *http.Request) {
	var req reportRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rp, err := s.svc.CreateReport(model.Report{
		ReporterID: req.ReporterID, TargetType: req.TargetType, TargetID: req.TargetID, Reason: req.Reason,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rp)
}

func (s *Server) listReports(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ReportFilter{
		TargetType: r.URL.Query().Get("target_type"),
		Status:     r.URL.Query().Get("status"),
		ReporterID: r.URL.Query().Get("reporter_id"),
	}
	items, total, err := s.svc.ListReports(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getReport(w http.ResponseWriter, r *http.Request) {
	rp, err := s.svc.GetReport(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rp)
}

func (s *Server) processReport(w http.ResponseWriter, r *http.Request) {
	var req reportActionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rp, err := s.svc.ProcessReport(r.PathValue("id"), req.Note)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rp)
}

func (s *Server) rejectReport(w http.ResponseWriter, r *http.Request) {
	var req reportActionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rp, err := s.svc.RejectReport(r.PathValue("id"), req.Note)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rp)
}

func (s *Server) deleteReport(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteReport(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
