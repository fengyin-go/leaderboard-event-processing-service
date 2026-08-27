package handler

import (
	"net/http"

	"leaderboard/internal/model"
	"leaderboard/pkg/httpx"
)

func (s *Server) registerScoreEventRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/score-events", s.createScoreEvent)
	mux.HandleFunc("GET /api/score-events", s.listScoreEvents)
	mux.HandleFunc("GET /api/score-events/{id}", s.getScoreEvent)
	mux.HandleFunc("DELETE /api/score-events/{id}", s.deleteScoreEvent)
}

type createScoreEventRequest struct {
	BoardID  string `json:"board_id"`
	MemberID string `json:"member_id"`
	SeasonID string `json:"season_id"`
	Value    int64  `json:"value"`
	Source   string `json:"source"`
}

func (s *Server) createScoreEvent(w http.ResponseWriter, r *http.Request) {
	var req createScoreEventRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.CreateScoreEvent(model.ScoreEvent{BoardID: req.BoardID, MemberID: req.MemberID, SeasonID: req.SeasonID, Value: req.Value, Source: req.Source})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, e)
}

func (s *Server) listScoreEvents(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ScoreEventFilter{
		BoardID:  r.URL.Query().Get("board_id"),
		MemberID: r.URL.Query().Get("member_id"),
		SeasonID: r.URL.Query().Get("season_id"),
		Source:   r.URL.Query().Get("source"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListScoreEvents(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getScoreEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, err := s.svc.GetScoreEvent(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) deleteScoreEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteScoreEvent(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
