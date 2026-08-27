package handler

import (
	"net/http"
	"strconv"

	"leaderboard/internal/model"
	"leaderboard/pkg/httpx"
)

func (s *Server) registerLeaderboardRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/leaderboards/{board_id}/overview", s.getLeaderboardOverview)
	mux.HandleFunc("GET /api/members/{member_id}/profile", s.getMemberProfile)
	mux.HandleFunc("GET /api/boards/{board_id}/season-comparison", s.compareSeasons)
	mux.HandleFunc("GET /api/boards/{board_id}/activity-trend", s.getActivityTrend)
	mux.HandleFunc("POST /api/bulk/members", s.bulkCreateMembers)
	mux.HandleFunc("POST /api/bulk/score-events", s.bulkCreateScoreEvents)
	mux.HandleFunc("POST /api/boards/{board_id}/recalculate", s.recalculateRanks)
	mux.HandleFunc("GET /api/stats", s.getSystemStats)
}

func (s *Server) getLeaderboardOverview(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	overview, err := s.svc.GetLeaderboardOverview(boardID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, overview)
}

func (s *Server) getMemberProfile(w http.ResponseWriter, r *http.Request) {
	memberID := r.PathValue("member_id")
	profile, err := s.svc.GetMemberLeaderboardProfile(memberID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, profile)
}

func (s *Server) compareSeasons(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	comparisons, err := s.svc.CompareSeasons(boardID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, comparisons)
}

func (s *Server) getActivityTrend(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days < 1 {
		days = 7
	}
	trend, err := s.svc.GetBoardActivityTrend(boardID, days)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, trend)
}

type bulkCreateMembersRequest struct {
	Members []model.Member `json:"members"`
}

func (s *Server) bulkCreateMembers(w http.ResponseWriter, r *http.Request) {
	var req bulkCreateMembersRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.BulkCreateMembers(req.Members)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, result)
}

type bulkCreateScoreEventsRequest struct {
	Events []model.ScoreEvent `json:"events"`
}

func (s *Server) bulkCreateScoreEvents(w http.ResponseWriter, r *http.Request) {
	var req bulkCreateScoreEventsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.BulkCreateScoreEvents(req.Events)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, result)
}

func (s *Server) recalculateRanks(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	seasonID := r.URL.Query().Get("season_id")
	if err := s.svc.RecalculateAllRanks(boardID, seasonID); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]string{"status": "recalculated"})
}

func (s *Server) getSystemStats(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetSystemStats()
	httpx.OK(w, stats)
}
