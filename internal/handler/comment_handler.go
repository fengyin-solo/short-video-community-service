package handler

import (
	"net/http"

	"shortvideo/internal/model"
	"shortvideo/pkg/httpx"
)

func (s *Server) registerCommentRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/comments", s.createComment)
	mux.HandleFunc("GET /api/comments", s.listComments)
	mux.HandleFunc("GET /api/comments/{id}", s.getComment)
	mux.HandleFunc("DELETE /api/comments/{id}", s.deleteComment)
}

type commentRequest struct {
	VideoID  string `json:"video_id"`
	UserID   string `json:"user_id"`
	ParentID string `json:"parent_id"`
	Content  string `json:"content"`
}

func (s *Server) createComment(w http.ResponseWriter, r *http.Request) {
	var req commentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateComment(model.Comment{
		VideoID: req.VideoID, UserID: req.UserID, ParentID: req.ParentID, Content: req.Content,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listComments(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CommentFilter{
		VideoID: r.URL.Query().Get("video_id"),
		UserID:  r.URL.Query().Get("user_id"),
		Status:  r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListComments(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getComment(w http.ResponseWriter, r *http.Request) {
	c, err := s.svc.GetComment(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteComment(w http.ResponseWriter, r *http.Request) {
	c, err := s.svc.DeleteComment(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}
