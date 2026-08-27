package model

import (
	"time"
)

type ChangeLog struct {
	ID       string    `json:"id"`
	BoardID  string    `json:"board_id"`
	MemberID string    `json:"member_id"`
	OldRank  int       `json:"old_rank"`
	NewRank  int       `json:"new_rank"`
	Delta    int       `json:"delta"`
	At       time.Time `json:"at"`
}

func (c *ChangeLog) Validate() error {
	if c.BoardID == "" {
		return NewValidationError("board_id", "所属榜单不能为空")
	}
	if c.MemberID == "" {
		return NewValidationError("member_id", "所属成员不能为空")
	}
	return nil
}

type ChangeLogFilter struct {
	BoardID  string
	MemberID string
}

func (f ChangeLogFilter) Match(c *ChangeLog) bool {
	if f.BoardID != "" && c.BoardID != f.BoardID {
		return false
	}
	if f.MemberID != "" && c.MemberID != f.MemberID {
		return false
	}
	return true
}
