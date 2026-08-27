package handler

import (
	"bytes"
	"io"
	"net/http"

	"leaderboard/pkg/httpx"
)

func (s *Server) registerImportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/import/members/json", s.importMembersJSON)
	mux.HandleFunc("POST /api/import/members/csv", s.importMembersCSV)
	mux.HandleFunc("POST /api/import/score-events/json", s.importScoreEventsJSON)
	mux.HandleFunc("GET /api/export/members/json", s.exportMembersJSON)
	mux.HandleFunc("GET /api/export/score-events/json", s.exportScoreEventsJSON)
	mux.HandleFunc("GET /api/export/rank-entries/csv", s.exportRankEntriesCSV)
}

func (s *Server) importMembersJSON(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(io.LimitReader(r.Body, 10<<20))
	if err != nil {
		httpx.BadRequest(w, "读取数据失败")
		return
	}
	if err := s.svc.ValidateImportData(data); err != nil {
		httpx.BadRequest(w, err.Error())
		return
	}
	result, err := s.svc.ImportMembersFromJSON(bytes.NewReader(data))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) importMembersCSV(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(io.LimitReader(r.Body, 10<<20))
	if err != nil {
		httpx.BadRequest(w, "读取数据失败")
		return
	}
	if err := s.svc.ValidateImportData(data); err != nil {
		httpx.BadRequest(w, err.Error())
		return
	}
	result, err := s.svc.ImportMembersFromCSV(bytes.NewReader(data))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) importScoreEventsJSON(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(io.LimitReader(r.Body, 10<<20))
	if err != nil {
		httpx.BadRequest(w, "读取数据失败")
		return
	}
	if err := s.svc.ValidateImportData(data); err != nil {
		httpx.BadRequest(w, err.Error())
		return
	}
	result, err := s.svc.ImportScoreEventsFromJSON(bytes.NewReader(data))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) exportMembersJSON(w http.ResponseWriter, r *http.Request) {
	data, err := s.svc.ExportMembersToJSON()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=members.json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (s *Server) exportScoreEventsJSON(w http.ResponseWriter, r *http.Request) {
	boardID := r.URL.Query().Get("board_id")
	data, err := s.svc.ExportScoreEventsToJSON(boardID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=score_events.json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (s *Server) exportRankEntriesCSV(w http.ResponseWriter, r *http.Request) {
	boardID := r.URL.Query().Get("board_id")
	seasonID := r.URL.Query().Get("season_id")
	data, err := s.svc.ExportRankEntriesToCSV(boardID, seasonID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=rank_entries.csv")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
