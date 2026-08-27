package handler

import (
	"net/http"
	"strconv"

	"leaderboard/internal/model"
	"leaderboard/pkg/httpx"
)

func (s *Server) registerRankEntryRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/rank-entries", s.listRankEntries)
	mux.HandleFunc("GET /api/rank-entries/{id}", s.getRankEntry)
	mux.HandleFunc("GET /api/boards/{board_id}/seasons/{season_id}/top", s.getTopN)
	mux.HandleFunc("GET /api/boards/{board_id}/seasons/{season_id}/range", s.getRankRange)
	mux.HandleFunc("GET /api/boards/{board_id}/seasons/{season_id}/members/{member_id}/rank", s.getMemberRank)
	mux.HandleFunc("GET /api/boards/{board_id}/seasons/{season_id}/report", s.getBoardReport)
}

func (s *Server) listRankEntries(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RankEntryFilter{
		BoardID:  r.URL.Query().Get("board_id"),
		SeasonID: r.URL.Query().Get("season_id"),
		MemberID: r.URL.Query().Get("member_id"),
	}
	items, total, err := s.svc.ListRankEntries(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRankEntry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	entry, err := s.svc.GetRankEntry(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, entry)
}

func (s *Server) getTopN(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	seasonID := r.PathValue("season_id")
	nStr := r.URL.Query().Get("n")
	n, _ := strconv.Atoi(nStr)
	if n < 1 {
		n = 10
	}
	items, err := s.svc.GetTopN(boardID, seasonID, n)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}

func (s *Server) getRankRange(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	seasonID := r.PathValue("season_id")
	start, _ := strconv.Atoi(r.URL.Query().Get("start"))
	end, _ := strconv.Atoi(r.URL.Query().Get("end"))
	if start < 1 {
		start = 1
	}
	if end < start {
		end = start + 9
	}
	items, err := s.svc.GetRankRange(boardID, seasonID, start, end)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}

func (s *Server) getMemberRank(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	seasonID := r.PathValue("season_id")
	memberID := r.PathValue("member_id")
	entry, err := s.svc.GetMemberRank(boardID, seasonID, memberID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, entry)
}

func (s *Server) getBoardReport(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	seasonID := r.PathValue("season_id")
	report, err := s.svc.BoardReport(boardID, seasonID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, report)
}
