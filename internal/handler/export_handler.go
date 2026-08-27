package handler

import (
	"net/http"

	"leaderboard/pkg/httpx"
)

func (s *Server) registerExportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/export/boards/{board_id}", s.exportBoard)
}

func (s *Server) exportBoard(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("board_id")
	data, err := s.svc.ExportBoardData(boardID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, data)
}
