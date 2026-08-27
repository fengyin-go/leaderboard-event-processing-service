package store

import (
	"sync"

	"leaderboard/internal/model"
)

type MemoryStore struct {
	mu            sync.RWMutex
	boards        map[string]*model.Board
	members       map[string]*model.Member
	seasons       map[string]*model.Season
	scoreEvents   map[string]*model.ScoreEvent
	rankEntries   map[string]*model.RankEntry
	rankSnapshots map[string]*model.RankSnapshot
	changeLogs    map[string]*model.ChangeLog
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		boards:        make(map[string]*model.Board),
		members:       make(map[string]*model.Member),
		seasons:       make(map[string]*model.Season),
		scoreEvents:   make(map[string]*model.ScoreEvent),
		rankEntries:   make(map[string]*model.RankEntry),
		rankSnapshots: make(map[string]*model.RankSnapshot),
		changeLogs:    make(map[string]*model.ChangeLog),
	}
}

var _ Store = (*MemoryStore)(nil)
