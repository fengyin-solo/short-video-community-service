package handler

import (
	"net/http"

	"shortvideo/internal/model"
	"shortvideo/pkg/httpx"
)

func (s *Server) registerLikeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/likes", s.likeVideo)
	mux.HandleFunc("DELETE /api/likes", s.unlikeVideo)
	mux.HandleFunc("GET /api/likes", s.listLikes)
}

type likeRequest struct {
	UserID  string `json:"user_id"`
	VideoID string `json:"video_id"`
}

func (s *Server) likeVideo(w http.ResponseWriter, r *http.Request) {
	var req likeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	l, err := s.svc.LikeVideo(req.UserID, req.VideoID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, l)
}

func (s *Server) unlikeVideo(w http.ResponseWriter, r *http.Request) {
	var req likeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.UnlikeVideo(req.UserID, req.VideoID); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) listLikes(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.LikeFilter{
		UserID:  r.URL.Query().Get("user_id"),
		VideoID: r.URL.Query().Get("video_id"),
	}
	items, total, err := s.svc.ListLikes(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}
