package handler

import (
	"net/http"

	"shortvideo/internal/model"
	"shortvideo/pkg/httpx"
)

func (s *Server) registerVideoTagRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/video-tags", s.attachTag)
	mux.HandleFunc("DELETE /api/video-tags", s.detachTag)
	mux.HandleFunc("GET /api/video-tags", s.listVideoTags)
	mux.HandleFunc("GET /api/videos/{id}/tags", s.listVideoTagNames)
}

type videoTagRequest struct {
	VideoID string `json:"video_id"`
	TagID   string `json:"tag_id"`
}

func (s *Server) attachTag(w http.ResponseWriter, r *http.Request) {
	var req videoTagRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	vt, err := s.svc.AttachTag(req.VideoID, req.TagID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, vt)
}

func (s *Server) detachTag(w http.ResponseWriter, r *http.Request) {
	var req videoTagRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.DetachTag(req.VideoID, req.TagID); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) listVideoTags(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.VideoTagFilter{
		VideoID: r.URL.Query().Get("video_id"),
		TagID:   r.URL.Query().Get("tag_id"),
	}
	items, total, err := s.svc.ListVideoTags(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) listVideoTagNames(w http.ResponseWriter, r *http.Request) {
	names := s.svc.ListVideoTagNames(r.PathValue("id"))
	httpx.OK(w, map[string]interface{}{"tags": names})
}
