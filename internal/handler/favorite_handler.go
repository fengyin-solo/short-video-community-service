package handler

import (
	"net/http"

	"shortvideo/internal/model"
	"shortvideo/pkg/httpx"
)

func (s *Server) registerFavoriteRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/favorites", s.addFavorite)
	mux.HandleFunc("DELETE /api/favorites", s.removeFavorite)
	mux.HandleFunc("GET /api/favorites", s.listFavorites)
}

type favoriteRequest struct {
	UserID  string `json:"user_id"`
	VideoID string `json:"video_id"`
	Folder  string `json:"folder"`
}

func (s *Server) addFavorite(w http.ResponseWriter, r *http.Request) {
	var req favoriteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.AddFavorite(req.UserID, req.VideoID, req.Folder)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, f)
}

func (s *Server) removeFavorite(w http.ResponseWriter, r *http.Request) {
	var req favoriteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.RemoveFavorite(req.UserID, req.VideoID); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) listFavorites(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.FavoriteFilter{
		UserID:  r.URL.Query().Get("user_id"),
		VideoID: r.URL.Query().Get("video_id"),
		Folder:  r.URL.Query().Get("folder"),
	}
	items, total, err := s.svc.ListFavorites(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}
