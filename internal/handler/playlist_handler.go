package handler

import (
	"net/http"

	"shortvideo/internal/model"
	"shortvideo/pkg/httpx"
)

func (s *Server) registerPlaylistRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/playlists", s.createPlaylist)
	mux.HandleFunc("GET /api/playlists", s.listPlaylists)
	mux.HandleFunc("GET /api/playlists/{id}", s.getPlaylist)
	mux.HandleFunc("PUT /api/playlists/{id}", s.updatePlaylist)
	mux.HandleFunc("DELETE /api/playlists/{id}", s.deletePlaylist)
	mux.HandleFunc("POST /api/playlists/{id}/videos", s.addVideoToPlaylist)
	mux.HandleFunc("DELETE /api/playlists/{id}/videos/{videoID}", s.removeVideoFromPlaylist)
}

type playlistRequest struct {
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (s *Server) createPlaylist(w http.ResponseWriter, r *http.Request) {
	var req playlistRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.CreatePlaylist(model.Playlist{
		UserID: req.UserID, Name: req.Name, Description: req.Description, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, p)
}

func (s *Server) listPlaylists(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.PlaylistFilter{
		UserID: r.URL.Query().Get("user_id"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListPlaylists(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getPlaylist(w http.ResponseWriter, r *http.Request) {
	p, err := s.svc.GetPlaylist(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) updatePlaylist(w http.ResponseWriter, r *http.Request) {
	var req playlistRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.UpdatePlaylist(r.PathValue("id"), model.Playlist{
		Name: req.Name, Description: req.Description, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) deletePlaylist(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeletePlaylist(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type playlistVideoRequest struct {
	VideoID string `json:"video_id"`
}

func (s *Server) addVideoToPlaylist(w http.ResponseWriter, r *http.Request) {
	var req playlistVideoRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.AddVideoToPlaylist(r.PathValue("id"), req.VideoID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) removeVideoFromPlaylist(w http.ResponseWriter, r *http.Request) {
	p, err := s.svc.RemoveVideoFromPlaylist(r.PathValue("id"), r.PathValue("videoID"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}
