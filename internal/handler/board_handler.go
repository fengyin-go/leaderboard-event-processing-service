package handler

import (
	"net/http"

	"leaderboard/internal/model"
	"leaderboard/pkg/httpx"
)

func (s *Server) registerBoardRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/boards", s.createBoard)
	mux.HandleFunc("GET /api/boards", s.listBoards)
	mux.HandleFunc("GET /api/boards/{id}", s.getBoard)
	mux.HandleFunc("PUT /api/boards/{id}", s.updateBoard)
	mux.HandleFunc("DELETE /api/boards/{id}", s.deleteBoard)
}

type createBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Metric      string `json:"metric"`
	SortOrder   string `json:"sort_order"`
	Status      string `json:"status"`
}

func (s *Server) createBoard(w http.ResponseWriter, r *http.Request) {
	var req createBoardRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.CreateBoard(model.Board{Name: req.Name, Description: req.Description, Metric: req.Metric, SortOrder: req.SortOrder, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, b)
}

func (s *Server) listBoards(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.BoardFilter{
		Metric:    r.URL.Query().Get("metric"),
		Status:    r.URL.Query().Get("status"),
		SortOrder: r.URL.Query().Get("sort_order"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListBoards(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getBoard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := s.svc.GetBoard(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

type updateBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Metric      string `json:"metric"`
	SortOrder   string `json:"sort_order"`
	Status      string `json:"status"`
}

func (s *Server) updateBoard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateBoardRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.UpdateBoard(id, model.Board{Name: req.Name, Description: req.Description, Metric: req.Metric, SortOrder: req.SortOrder, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

func (s *Server) deleteBoard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteBoard(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
