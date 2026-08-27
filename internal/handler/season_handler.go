package handler

import (
	"net/http"
	"time"

	"leaderboard/internal/model"
	"leaderboard/pkg/httpx"
)

func (s *Server) registerSeasonRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/seasons", s.createSeason)
	mux.HandleFunc("GET /api/seasons", s.listSeasons)
	mux.HandleFunc("GET /api/seasons/{id}", s.getSeason)
	mux.HandleFunc("PUT /api/seasons/{id}", s.updateSeason)
	mux.HandleFunc("POST /api/seasons/{id}/transition", s.transitionSeason)
	mux.HandleFunc("DELETE /api/seasons/{id}", s.deleteSeason)
}

type createSeasonRequest struct {
	BoardID string    `json:"board_id"`
	Name    string    `json:"name"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
	Status  string    `json:"status"`
}

func (s *Server) createSeason(w http.ResponseWriter, r *http.Request) {
	var req createSeasonRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	se, err := s.svc.CreateSeason(model.Season{BoardID: req.BoardID, Name: req.Name, StartAt: req.StartAt, EndAt: req.EndAt, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, se)
}

func (s *Server) listSeasons(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SeasonFilter{
		BoardID: r.URL.Query().Get("board_id"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListSeasons(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSeason(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	se, err := s.svc.GetSeason(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, se)
}

type updateSeasonRequest struct {
	Name    string    `json:"name"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}

func (s *Server) updateSeason(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateSeasonRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	se, err := s.svc.UpdateSeason(id, model.Season{Name: req.Name, StartAt: req.StartAt, EndAt: req.EndAt})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, se)
}

type transitionSeasonRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionSeason(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionSeasonRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	se, err := s.svc.TransitionSeasonStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, se)
}

func (s *Server) deleteSeason(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSeason(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
