package handler

import (
	"net/http"

	"shortvideo/internal/model"
	"shortvideo/pkg/httpx"
)

func (s *Server) registerUserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/users", s.createUser)
	mux.HandleFunc("GET /api/users", s.listUsers)
	mux.HandleFunc("GET /api/users/{id}", s.getUser)
	mux.HandleFunc("PUT /api/users/{id}", s.updateUser)
	mux.HandleFunc("DELETE /api/users/{id}", s.deleteUser)
	mux.HandleFunc("POST /api/users/{id}/ban", s.banUser)
	mux.HandleFunc("POST /api/users/{id}/unban", s.unbanUser)
}

type userRequest struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Status   string `json:"status"`
}

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	u, err := s.svc.CreateUser(model.User{
		Username: req.Username, Nickname: req.Nickname, Avatar: req.Avatar, Bio: req.Bio, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, u)
}

func (s *Server) listUsers(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.UserFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListUsers(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	u, err := s.svc.GetUser(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, u)
}

func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	u, err := s.svc.UpdateUser(r.PathValue("id"), model.User{
		Username: req.Username, Nickname: req.Nickname, Avatar: req.Avatar, Bio: req.Bio, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, u)
}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteUser(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) banUser(w http.ResponseWriter, r *http.Request) {
	u, err := s.svc.BanUser(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, u)
}

func (s *Server) unbanUser(w http.ResponseWriter, r *http.Request) {
	u, err := s.svc.UnbanUser(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, u)
}
