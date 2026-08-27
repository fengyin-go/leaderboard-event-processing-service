package handler

import (
	"net/http"
	"strconv"

	"leaderboard/pkg/httpx"
)

func (s *Server) registerEventExtRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/score-events/{id}/revert", s.revertScoreEvent)
	mux.HandleFunc("GET /api/boards/{board_id}/score-stats", s.getScoreEventStats)
	mux.HandleFunc("GET /api/boards/{board_id}/members/{member_id}/timeline", s.getScoreEventTimeline)
	mux.HandleFunc("POST /api/rank-snapshots/compare", s.compareSnapshots)
	mux.HandleFunc("GET /api/boards/{board_id}/change-log-stats", s.getChangeLogStats)
}

func (s *Server) revertScoreEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, err := s.svc.RevertScoreEvent(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, e)
}

func (s *Server) getScoreEventStats(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	seasonID := r.URL.Query().Get("season_id")
	stats, err := s.svc.GetScoreEventStats(boardID, seasonID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) getScoreEventTimeline(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	memberID := r.PathValue("member_id")
	limit := 50
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	timeline, err := s.svc.GetScoreEventTimeline(boardID, memberID, limit)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, timeline)
}

type compareSnapshotsRequest struct {
	SnapshotID1 string `json:"snapshot_id_1"`
	SnapshotID2 string `json:"snapshot_id_2"`
}

func (s *Server) compareSnapshots(w http.ResponseWriter, r *http.Request) {
	var req compareSnapshotsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.CompareSnapshots(req.SnapshotID1, req.SnapshotID2)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) getChangeLogStats(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	stats, err := s.svc.GetChangeLogStats(boardID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}
