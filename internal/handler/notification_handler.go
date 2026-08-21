package handler

import (
	"net/http"
	"strconv"

	"shortvideo/internal/model"
	"shortvideo/pkg/httpx"
)

func (s *Server) registerNotificationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/notifications", s.createNotification)
	mux.HandleFunc("GET /api/notifications", s.listNotifications)
	mux.HandleFunc("GET /api/notifications/{id}", s.getNotification)
	mux.HandleFunc("POST /api/notifications/{id}/read", s.markNotificationRead)
	mux.HandleFunc("POST /api/notifications/read-all", s.markAllNotificationsRead)
	mux.HandleFunc("DELETE /api/notifications/{id}", s.deleteNotification)
}

type notificationRequest struct {
	UserID  string `json:"user_id"`
	Type    string `json:"type"`
	RefID   string `json:"ref_id"`
	Content string `json:"content"`
}

func (s *Server) createNotification(w http.ResponseWriter, r *http.Request) {
	var req notificationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.CreateNotification(model.Notification{
		UserID: req.UserID, Type: req.Type, RefID: req.RefID, Content: req.Content,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, n)
}

func (s *Server) listNotifications(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.NotificationFilter{
		UserID: r.URL.Query().Get("user_id"),
		Type:   r.URL.Query().Get("type"),
	}
	if v := r.URL.Query().Get("read"); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			filter.Read = &b
		}
	}
	items, total, err := s.svc.ListNotifications(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getNotification(w http.ResponseWriter, r *http.Request) {
	n, err := s.svc.GetNotification(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

func (s *Server) markNotificationRead(w http.ResponseWriter, r *http.Request) {
	n, err := s.svc.MarkNotificationRead(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

type readAllRequest struct {
	UserID string `json:"user_id"`
}

func (s *Server) markAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	var req readAllRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	count, err := s.svc.MarkAllNotificationsRead(req.UserID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"marked": count})
}

func (s *Server) deleteNotification(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteNotification(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
