package handler

import (
	"net/http"

	"leaderboard/pkg/httpx"
)

func (s *Server) registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", s.healthCheck)
	mux.HandleFunc("GET /ready", s.readyCheck)
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, map[string]string{"status": "healthy"})
}

func (s *Server) readyCheck(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, map[string]string{"status": "ready"})
}
