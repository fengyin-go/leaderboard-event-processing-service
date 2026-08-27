package handler

import (
	"net/http"

	"leaderboard/internal/model"
	"leaderboard/pkg/httpx"
)

func (s *Server) registerMemberRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/members", s.createMember)
	mux.HandleFunc("GET /api/members", s.listMembers)
	mux.HandleFunc("GET /api/members/{id}", s.getMember)
	mux.HandleFunc("PUT /api/members/{id}", s.updateMember)
	mux.HandleFunc("DELETE /api/members/{id}", s.deleteMember)
}

type createMemberRequest struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

func (s *Server) createMember(w http.ResponseWriter, r *http.Request) {
	var req createMemberRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.CreateMember(model.Member{Name: req.Name, Tag: req.Tag})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, m)
}

func (s *Server) listMembers(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MemberFilter{
		Tag:     r.URL.Query().Get("tag"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListMembers(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMember(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	m, err := s.svc.GetMember(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

type updateMemberRequest struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

func (s *Server) updateMember(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateMemberRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.UpdateMember(id, model.Member{Name: req.Name, Tag: req.Tag})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) deleteMember(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteMember(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
