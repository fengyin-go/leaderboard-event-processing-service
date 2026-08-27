// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"leaderboard/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Board
	CreateBoard(b *model.Board) error
	GetBoard(id string) (*model.Board, error)
	GetBoardByName(name string) (*model.Board, error)
	ListBoards() []*model.Board
	UpdateBoard(b *model.Board) error
	DeleteBoard(id string) error

	// Member
	CreateMember(m *model.Member) error
	GetMember(id string) (*model.Member, error)
	GetMemberByTag(tag string) (*model.Member, error)
	ListMembers() []*model.Member
	UpdateMember(m *model.Member) error
	DeleteMember(id string) error

	// Season
	CreateSeason(s *model.Season) error
	GetSeason(id string) (*model.Season, error)
	ListSeasons() []*model.Season
	UpdateSeason(s *model.Season) error
	DeleteSeason(id string) error

	// ScoreEvent
	CreateScoreEvent(e *model.ScoreEvent) error
	GetScoreEvent(id string) (*model.ScoreEvent, error)
	ListScoreEvents() []*model.ScoreEvent
	DeleteScoreEvent(id string) error

	// RankEntry
	CreateRankEntry(r *model.RankEntry) error
	GetRankEntry(id string) (*model.RankEntry, error)
	GetRankEntryByBoardSeasonMember(boardID, seasonID, memberID string) (*model.RankEntry, error)
	ListRankEntries() []*model.RankEntry
	UpdateRankEntry(r *model.RankEntry) error
	DeleteRankEntry(id string) error

	// RankSnapshot
	CreateRankSnapshot(rs *model.RankSnapshot) error
	GetRankSnapshot(id string) (*model.RankSnapshot, error)
	ListRankSnapshots() []*model.RankSnapshot
	DeleteRankSnapshot(id string) error

	// ChangeLog
	CreateChangeLog(c *model.ChangeLog) error
	GetChangeLog(id string) (*model.ChangeLog, error)
	ListChangeLogs() []*model.ChangeLog
	DeleteChangeLog(id string) error
}
