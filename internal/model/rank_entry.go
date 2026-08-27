package model

import (
	"time"
)

type RankEntry struct {
	ID       string    `json:"id"`
	BoardID  string    `json:"board_id"`
	SeasonID string    `json:"season_id"`
	MemberID string    `json:"member_id"`
	Score    int64     `json:"score"`
	Rank     int       `json:"rank"`
	PrevRank int       `json:"prev_rank"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *RankEntry) Validate() error {
	if r.BoardID == "" {
		return NewValidationError("board_id", "所属榜单不能为空")
	}
	if r.MemberID == "" {
		return NewValidationError("member_id", "所属成员不能为空")
	}
	return nil
}

type RankEntryFilter struct {
	BoardID  string
	SeasonID string
	MemberID string
}

func (f RankEntryFilter) Match(r *RankEntry) bool {
	if f.BoardID != "" && r.BoardID != f.BoardID {
		return false
	}
	if f.SeasonID != "" && r.SeasonID != f.SeasonID {
		return false
	}
	if f.MemberID != "" && r.MemberID != f.MemberID {
		return false
	}
	return true
}
