package handler

import (
	"net/http"

	"leaderboard/internal/model"
	"leaderboard/pkg/httpx"
)

func (s *Server) registerChangeLogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/change-logs", s.listChangeLogs)
	mux.HandleFunc("GET /api/change-logs/{id}", s.getChangeLog)
	mux.HandleFunc("DELETE /api/change-logs/{id}", s.deleteChangeLog)
}

func (s *Server) listChangeLogs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ChangeLogFilter{
		BoardID:  r.URL.Query().Get("board_id"),
		MemberID: r.URL.Query().Get("member_id"),
	}
	items, total, err := s.svc.ListChangeLogs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getChangeLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.GetChangeLog(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteChangeLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteChangeLog(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
