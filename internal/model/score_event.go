package model

import (
	"strings"
	"time"
)

type ScoreEvent struct {
	ID       string    `json:"id"`
	BoardID  string    `json:"board_id"`
	MemberID string    `json:"member_id"`
	SeasonID string    `json:"season_id"`
	Value    int64     `json:"value"`
	Source   string    `json:"source"`
	At       time.Time `json:"at"`
}

func (e *ScoreEvent) Validate() error {
	e.Source = strings.TrimSpace(e.Source)

	if e.BoardID == "" {
		return NewValidationError("board_id", "所属榜单不能为空")
	}
	if e.MemberID == "" {
		return NewValidationError("member_id", "所属成员不能为空")
	}
	if e.Source == "" {
		return NewValidationError("source", "事件来源不能为空")
	}
	return nil
}

type ScoreEventFilter struct {
	BoardID  string
	MemberID string
	SeasonID string
	Source   string
	Keyword  string
}

func (f ScoreEventFilter) Match(e *ScoreEvent) bool {
	if f.BoardID != "" && e.BoardID != f.BoardID {
		return false
	}
	if f.MemberID != "" && e.MemberID != f.MemberID {
		return false
	}
	if f.SeasonID != "" && e.SeasonID != f.SeasonID {
		return false
	}
	if f.Source != "" && e.Source != f.Source {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(e.Source), k) {
			return false
		}
	}
	return true
}
