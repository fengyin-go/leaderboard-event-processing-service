package handler

import (
	"net/http"

	"leaderboard/internal/model"
	"leaderboard/pkg/httpx"
)

func (s *Server) registerRankSnapshotRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/rank-snapshots", s.createRankSnapshot)
	mux.HandleFunc("GET /api/rank-snapshots", s.listRankSnapshots)
	mux.HandleFunc("GET /api/rank-snapshots/{id}", s.getRankSnapshot)
	mux.HandleFunc("DELETE /api/rank-snapshots/{id}", s.deleteRankSnapshot)
}

type createRankSnapshotRequest struct {
	BoardID  string `json:"board_id"`
	SeasonID string `json:"season_id"`
}

func (s *Server) createRankSnapshot(w http.ResponseWriter, r *http.Request) {
	var req createRankSnapshotRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rs, err := s.svc.CaptureRankSnapshot(req.BoardID, req.SeasonID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rs)
}

func (s *Server) listRankSnapshots(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RankSnapshotFilter{
		BoardID:  r.URL.Query().Get("board_id"),
		SeasonID: r.URL.Query().Get("season_id"),
	}
	items, total, err := s.svc.ListRankSnapshots(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRankSnapshot(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rs, err := s.svc.GetRankSnapshot(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rs)
}

func (s *Server) deleteRankSnapshot(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRankSnapshot(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
