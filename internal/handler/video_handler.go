package handler

import (
	"net/http"

	"shortvideo/internal/model"
	"shortvideo/pkg/httpx"
)

func (s *Server) registerVideoRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/videos", s.createVideo)
	mux.HandleFunc("GET /api/videos", s.listVideos)
	mux.HandleFunc("GET /api/videos/{id}", s.getVideo)
	mux.HandleFunc("DELETE /api/videos/{id}", s.deleteVideo)
	mux.HandleFunc("POST /api/videos/{id}/submit", s.submitVideo)
	mux.HandleFunc("POST /api/videos/{id}/approve", s.approveVideo)
	mux.HandleFunc("POST /api/videos/{id}/ban", s.banVideo)
	mux.HandleFunc("POST /api/videos/{id}/view", s.viewVideo)
}

type videoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	VideoURL    string `json:"video_url"`
	AuthorID    string `json:"author_id"`
	CategoryID  string `json:"category_id"`
	Duration    int    `json:"duration"`
}

func (s *Server) createVideo(w http.ResponseWriter, r *http.Request) {
	var req videoRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	v, err := s.svc.CreateVideo(model.Video{
		Title: req.Title, Description: req.Description, CoverURL: req.CoverURL, VideoURL: req.VideoURL,
		AuthorID: req.AuthorID, CategoryID: req.CategoryID, Duration: req.Duration,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, v)
}

func (s *Server) listVideos(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.VideoFilter{
		Status:     r.URL.Query().Get("status"),
		CategoryID: r.URL.Query().Get("category_id"),
		AuthorID:   r.URL.Query().Get("author_id"),
		Keyword:    r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListVideos(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getVideo(w http.ResponseWriter, r *http.Request) {
	v, err := s.svc.GetVideo(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}

func (s *Server) deleteVideo(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteVideo(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) submitVideo(w http.ResponseWriter, r *http.Request) {
	v, err := s.svc.SubmitVideoForReview(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}

func (s *Server) approveVideo(w http.ResponseWriter, r *http.Request) {
	v, err := s.svc.ApproveVideo(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}

func (s *Server) banVideo(w http.ResponseWriter, r *http.Request) {
	v, err := s.svc.BanVideo(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}

func (s *Server) viewVideo(w http.ResponseWriter, r *http.Request) {
	v, err := s.svc.IncrementView(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}
