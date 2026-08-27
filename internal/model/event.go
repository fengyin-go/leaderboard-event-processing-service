package model

import (
	"strings"
	"time"
)

const (
	EventTypeScoreEvent = "score_event"
	EventTypeRankChange = "rank_change"
	EventTypeSeasonTransition = "season_transition"
	EventTypeSnapshot = "snapshot"
)

const (
	EventLevelInfo = "info"
	EventLevelWarn = "warn"
	EventLevelError = "error"
)

type Event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Level     string    `json:"level"`
	BoardID   string    `json:"board_id,omitempty"`
	MemberID  string    `json:"member_id,omitempty"`
	SeasonID  string    `json:"season_id,omitempty"`
	Message   string    `json:"message"`
	Payload   string    `json:"payload,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (e *Event) Validate() error {
	e.Type = strings.TrimSpace(e.Type)
	e.Level = strings.TrimSpace(e.Level)
	e.Message = strings.TrimSpace(e.Message)
	if e.Type == "" {
		return NewValidationError("type", "事件类型不能为空")
	}
	if e.Message == "" {
		return NewValidationError("message", "事件消息不能为空")
	}
	if e.Level == "" {
		e.Level = EventLevelInfo
	}
	if e.Level != EventLevelInfo && e.Level != EventLevelWarn && e.Level != EventLevelError {
		return NewValidationError("level", "事件级别不合法")
	}
	return nil
}

type EventFilter struct {
	Type     string
	Level    string
	BoardID  string
	MemberID string
	SeasonID string
	Keyword  string
}

func (f EventFilter) Match(e *Event) bool {
	if f.Type != "" && e.Type != f.Type {
		return false
	}
	if f.Level != "" && e.Level != f.Level {
		return false
	}
	if f.BoardID != "" && e.BoardID != f.BoardID {
		return false
	}
	if f.MemberID != "" && e.MemberID != f.MemberID {
		return false
	}
	if f.SeasonID != "" && e.SeasonID != f.SeasonID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(e.Message), k) {
			return false
		}
	}
	return true
}
