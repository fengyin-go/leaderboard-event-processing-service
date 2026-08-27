package model

import (
	"strings"
	"time"
)

const (
	SeasonStatusUpcoming = "upcoming"
	SeasonStatusActive   = "active"
	SeasonStatusFinished = "finished"
)

var seasonTransitions = map[string]map[string]bool{
	SeasonStatusUpcoming: {SeasonStatusActive: true},
	SeasonStatusActive:   {SeasonStatusFinished: true},
	SeasonStatusFinished: {},
}

func SeasonCanTransition(from, to string) bool {
	if m, ok := seasonTransitions[from]; ok {
		return m[to]
	}
	return false
}

type Season struct {
	ID        string    `json:"id"`
	BoardID   string    `json:"board_id"`
	Name      string    `json:"name"`
	StartAt   time.Time `json:"start_at"`
	EndAt     time.Time `json:"end_at"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Season) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	s.Status = strings.TrimSpace(s.Status)

	if s.Name == "" {
		return NewValidationError("name", "赛季名称不能为空")
	}
	if s.BoardID == "" {
		return NewValidationError("board_id", "所属榜单不能为空")
	}
	if s.Status == "" {
		s.Status = SeasonStatusUpcoming
	}
	if s.Status != SeasonStatusUpcoming && s.Status != SeasonStatusActive && s.Status != SeasonStatusFinished {
		return NewValidationError("status", "赛季状态不合法")
	}
	if !s.EndAt.IsZero() && !s.StartAt.IsZero() && s.EndAt.Before(s.StartAt) {
		return NewValidationError("end_at", "结束时间不能早于开始时间")
	}
	return nil
}

type SeasonFilter struct {
	BoardID string
	Status  string
	Name    string
	Keyword string
}

func (f SeasonFilter) Match(s *Season) bool {
	if f.BoardID != "" && s.BoardID != f.BoardID {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.Name != "" && s.Name != f.Name {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.Name), k) {
			return false
		}
	}
	return true
}
