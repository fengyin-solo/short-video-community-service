package handler

import (
	"net/http"

	"shortvideo/internal/model"
	"shortvideo/pkg/httpx"
)

func (s *Server) registerFollowRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/follows", s.followUser)
	mux.HandleFunc("DELETE /api/follows", s.unfollowUser)
	mux.HandleFunc("GET /api/follows", s.listFollows)
}

type followRequest struct {
	FollowerID string `json:"follower_id"`
	FolloweeID string `json:"followee_id"`
}

func (s *Server) followUser(w http.ResponseWriter, r *http.Request) {
	var req followRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.FollowUser(req.FollowerID, req.FolloweeID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, f)
}

func (s *Server) unfollowUser(w http.ResponseWriter, r *http.Request) {
	var req followRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.UnfollowUser(req.FollowerID, req.FolloweeID); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) listFollows(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.FollowFilter{
		FollowerID: r.URL.Query().Get("follower_id"),
		FolloweeID: r.URL.Query().Get("followee_id"),
	}
	items, total, err := s.svc.ListFollows(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}
