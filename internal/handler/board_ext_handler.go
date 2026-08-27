package handler

import (
	"net/http"

	"leaderboard/pkg/httpx"
)

func (s *Server) registerBoardExtRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/boards/{id}/archive", s.archiveBoard)
	mux.HandleFunc("POST /api/boards/{id}/activate", s.activateBoard)
	mux.HandleFunc("POST /api/boards/{id}/duplicate", s.duplicateBoard)
	mux.HandleFunc("GET /api/boards/{id}/seasons", s.getBoardSeasons)
	mux.HandleFunc("GET /api/boards/{id}/stats", s.getBoardStats)
}

func (s *Server) archiveBoard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := s.svc.ArchiveBoard(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

func (s *Server) activateBoard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := s.svc.ActivateBoard(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

type duplicateBoardRequest struct {
	Name string `json:"name"`
}

func (s *Server) duplicateBoard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req duplicateBoardRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.DuplicateBoard(id, req.Name)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, b)
}

func (s *Server) getBoardSeasons(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	seasons, err := s.svc.GetBoardSeasons(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, seasons)
}

func (s *Server) getBoardStats(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	stats, err := s.svc.GetBoardStats(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}
