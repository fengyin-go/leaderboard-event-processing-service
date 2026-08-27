package model

import (
	"time"
)

type RankSnapshot struct {
	ID          string    `json:"id"`
	BoardID     string    `json:"board_id"`
	SeasonID    string    `json:"season_id"`
	CapturedAt  time.Time `json:"captured_at"`
	EntriesJSON string    `json:"entries_json"`
	CreatedAt   time.Time `json:"created_at"`
}

func (rs *RankSnapshot) Validate() error {
	if rs.BoardID == "" {
		return NewValidationError("board_id", "所属榜单不能为空")
	}
	if rs.EntriesJSON == "" {
		return NewValidationError("entries_json", "快照内容不能为空")
	}
	return nil
}

type RankSnapshotFilter struct {
	BoardID  string
	SeasonID string
}

func (f RankSnapshotFilter) Match(rs *RankSnapshot) bool {
	if f.BoardID != "" && rs.BoardID != f.BoardID {
		return false
	}
	if f.SeasonID != "" && rs.SeasonID != f.SeasonID {
		return false
	}
	return true
}
