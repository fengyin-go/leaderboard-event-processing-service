package handler

import (
	"net/http"
	"time"

	"leaderboard/pkg/httpx"
)

func (s *Server) registerSeasonExtRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/seasons/{id}/finalize", s.finalizeSeason)
	mux.HandleFunc("POST /api/seasons/{id}/clone", s.cloneSeason)
	mux.HandleFunc("GET /api/boards/{board_id}/current-season", s.getCurrentSeason)
}

func (s *Server) finalizeSeason(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	se, err := s.svc.FinalizeSeason(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, se)
}

type cloneSeasonRequest struct {
	Name    string    `json:"name"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}

func (s *Server) cloneSeason(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req cloneSeasonRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	se, err := s.svc.CloneSeason(id, req.Name, req.StartAt, req.EndAt)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, se)
}

func (s *Server) getCurrentSeason(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	se, err := s.svc.GetCurrentSeason(boardID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, se)
}
