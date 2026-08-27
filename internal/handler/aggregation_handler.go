package handler

import (
	"net/http"
	"strconv"
	"time"

	"leaderboard/pkg/httpx"
)

func (s *Server) registerAggregationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/boards/{board_id}/member-activity", s.getMemberActivity)
	mux.HandleFunc("GET /api/boards/{board_id}/score-distribution", s.getScoreDistribution)
	mux.HandleFunc("GET /api/boards/{board_id}/hot-streaks", s.getHotStreaks)
	mux.HandleFunc("GET /api/boards/{board_id}/milestones", s.getSeasonMilestones)
	mux.HandleFunc("GET /api/boards/{board_id}/rank-velocity", s.getRankVelocity)
}

func (s *Server) getMemberActivity(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	items, err := s.svc.GetMemberActivity(boardID, limit)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}

func (s *Server) getScoreDistribution(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	seasonID := r.URL.Query().Get("season_id")
	buckets, _ := strconv.Atoi(r.URL.Query().Get("buckets"))
	items, err := s.svc.GetScoreDistribution(boardID, seasonID, buckets)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}

func (s *Server) getHotStreaks(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	items, err := s.svc.GetHotStreaks(boardID, limit)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}

func (s *Server) getSeasonMilestones(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	items, err := s.svc.GetSeasonMilestones(boardID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}

func (s *Server) getRankVelocity(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	seasonID := r.URL.Query().Get("season_id")
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days < 1 {
		days = 7
	}
	since := time.Now().AddDate(0, 0, -days)
	items, err := s.svc.GetRankVelocity(boardID, seasonID, since)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}
