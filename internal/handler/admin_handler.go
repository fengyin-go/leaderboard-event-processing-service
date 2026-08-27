package handler

import (
	"net/http"
	"strconv"
	"time"

	"leaderboard/pkg/httpx"
)

func (s *Server) registerAdminRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/admin/boards", s.adminListBoards)
	mux.HandleFunc("DELETE /api/admin/boards/{board_id}", s.adminPurgeBoard)
	mux.HandleFunc("GET /api/admin/report", s.adminGlobalReport)
	mux.HandleFunc("GET /api/admin/inactive-members", s.adminInactiveMembers)
	mux.HandleFunc("POST /api/admin/boards/{board_id}/rebuild-logs", s.adminRebuildLogs)
}

func (s *Server) adminListBoards(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	items, total, err := s.svc.AdminListAllBoards(pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) adminPurgeBoard(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	if err := s.svc.AdminPurgeBoardData(boardID); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) adminGlobalReport(w http.ResponseWriter, r *http.Request) {
	report, err := s.svc.AdminGlobalReport()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, report)
}

func (s *Server) adminInactiveMembers(w http.ResponseWriter, r *http.Request) {
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days < 1 {
		days = 30
	}
	since := time.Now().AddDate(0, 0, -days)
	members, err := s.svc.AdminListInactiveMembers(since)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, members)
}

func (s *Server) adminRebuildLogs(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	if err := s.svc.AdminRebuildAllChangeLogs(boardID); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]string{"status": "rebuilt"})
}
