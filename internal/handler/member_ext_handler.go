package handler

import (
	"net/http"
	"strconv"

	"leaderboard/pkg/httpx"
)

func (s *Server) registerMemberExtRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/member-lookup/{tag}", s.getMemberByTag)
	mux.HandleFunc("GET /api/members/{id}/score-events", s.getMemberScoreEvents)
	mux.HandleFunc("GET /api/members/{id}/total-score", s.getMemberTotalScore)
	mux.HandleFunc("GET /api/members/{id}/rank-history", s.getMemberRankHistory)
}

func (s *Server) getMemberByTag(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")
	m, err := s.svc.SearchMembersByTag(tag)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) getMemberScoreEvents(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	limit := 20
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	events, err := s.svc.GetMemberScoreEvents(id, limit)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, events)
}

func (s *Server) getMemberTotalScore(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	score := s.svc.GetMemberTotalScore(id)
	httpx.OK(w, map[string]interface{}{"member_id": id, "total_score": score})
}

func (s *Server) getMemberRankHistory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	boardID := r.URL.Query().Get("board_id")
	history, err := s.svc.GetMemberRankHistory(id, boardID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, history)
}
