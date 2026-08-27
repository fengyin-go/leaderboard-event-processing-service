// Package handler 实现 HTTP 处理器层。
package handler

import (
	"errors"
	"net/http"
	"time"

	"leaderboard/internal/config"
	"leaderboard/internal/middleware"
	"leaderboard/internal/model"
	"leaderboard/internal/service"
	"leaderboard/internal/store"
	"leaderboard/pkg/httpx"
	"leaderboard/pkg/logger"
)

type Server struct {
	svc *service.Service
	log *logger.Logger
	cfg *config.Config
}

func NewServer(svc *service.Service, log *logger.Logger, cfg *config.Config) *Server {
	return &Server{svc: svc, log: log, cfg: cfg}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	s.registerBoardRoutes(mux)
	s.registerMemberRoutes(mux)
	s.registerSeasonRoutes(mux)
	s.registerScoreEventRoutes(mux)
	s.registerRankEntryRoutes(mux)
	s.registerRankSnapshotRoutes(mux)
	s.registerChangeLogRoutes(mux)
	s.registerExportRoutes(mux)
	s.registerLeaderboardRoutes(mux)
	s.registerHealthRoutes(mux)
	s.registerImportRoutes(mux)
	s.registerAggregationRoutes(mux)
	s.registerBoardExtRoutes(mux)
	s.registerMemberExtRoutes(mux)
	s.registerSeasonExtRoutes(mux)
	s.registerEventExtRoutes(mux)
	s.registerAdminRoutes(mux)

	var h http.Handler = mux
	h = middleware.RecoveryMiddleware(s.log)(h)
	h = middleware.LoggingMiddleware(s.log)(h)
	h = middleware.AuthMiddleware(s.cfg)(h)
	h = middleware.CORSMiddleware(h)
	h = middleware.RequestIDMiddleware(h)
	h = middleware.RateLimitMiddleware(middleware.NewRateLimiter(100, time.Minute))(h)
	return h
}

func (s *Server) maxPageSize() int {
	if s.cfg != nil && s.cfg.MaxPageSize > 0 {
		return s.cfg.MaxPageSize
	}
	return 100
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case model.IsValidationError(err):
		httpx.BadRequest(w, err.Error())
	case errors.Is(err, store.ErrNotFound):
		httpx.NotFound(w, err.Error())
	case errors.Is(err, store.ErrConflict):
		httpx.Conflict(w, err.Error())
	default:
		httpx.InternalError(w, err.Error())
	}
}
