package service

import (
	"encoding/json"
	"sort"
	"time"

	"leaderboard/internal/model"
)

func (s *Service) RevertScoreEvent(id string) (*model.ScoreEvent, error) {
	e, err := s.store.GetScoreEvent(id)
	if err != nil {
		return nil, err
	}
	e.Value = -e.Value
	e.Source = "revert:" + e.Source
	e.At = time.Now()
	reverted, err := s.CreateScoreEvent(*e)
	if err != nil {
		return nil, err
	}
	return reverted, nil
}

func (s *Service) GetScoreEventStats(boardID, seasonID string) (map[string]interface{}, error) {
	all := s.store.ListScoreEvents()
	var total int64
	count := 0
	sources := make(map[string]int)
	for _, e := range all {
		if e.BoardID == boardID && (seasonID == "" || e.SeasonID == seasonID) {
			total += e.Value
			count++
			sources[e.Source]++
		}
	}
	return map[string]interface{}{
		"board_id":    boardID,
		"season_id":   seasonID,
		"total_score": total,
		"event_count": count,
		"sources":     sources,
	}, nil
}

func (s *Service) GetScoreEventTimeline(boardID, memberID string, limit int) ([]*model.ScoreEvent, error) {
	if limit < 1 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	all := s.store.ListScoreEvents()
	var result []*model.ScoreEvent
	for _, e := range all {
		if e.BoardID == boardID && e.MemberID == memberID {
			result = append(result, e)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].At.After(result[j].At)
	})
	if len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

func (s *Service) CompareSnapshots(snapshotID1, snapshotID2 string) (map[string]interface{}, error) {
	s1, err := s.store.GetRankSnapshot(snapshotID1)
	if err != nil {
		return nil, err
	}
	s2, err := s.store.GetRankSnapshot(snapshotID2)
	if err != nil {
		return nil, err
	}
	var entries1, entries2 []*model.RankEntry
	_ = json.Unmarshal([]byte(s1.EntriesJSON), &entries1)
	_ = json.Unmarshal([]byte(s2.EntriesJSON), &entries2)

	m1 := make(map[string]*model.RankEntry)
	for _, e := range entries1 {
		m1[e.MemberID] = e
	}
	m2 := make(map[string]*model.RankEntry)
	for _, e := range entries2 {
		m2[e.MemberID] = e
	}

	type change struct {
		MemberID string `json:"member_id"`
		OldRank  int    `json:"old_rank"`
		NewRank  int    `json:"new_rank"`
		Delta    int    `json:"delta"`
	}
	var changes []change
	for memberID, e2 := range m2 {
		if e1, ok := m1[memberID]; ok {
			if e1.Rank != e2.Rank {
				changes = append(changes, change{
					MemberID: memberID,
					OldRank:  e1.Rank,
					NewRank:  e2.Rank,
					Delta:    e1.Rank - e2.Rank,
				})
			}
		}
	}
	return map[string]interface{}{
		"snapshot_1": s1.ID,
		"snapshot_2": s2.ID,
		"changes":    changes,
		"compared_at": time.Now(),
	}, nil
}

func (s *Service) GetChangeLogStats(boardID string) (map[string]interface{}, error) {
	logs := s.store.ListChangeLogs()
	var totalMoves int
	upMoves := 0
	downMoves := 0
	memberMoves := make(map[string]int)
	for _, c := range logs {
		if c.BoardID != boardID {
			continue
		}
		totalMoves++
		if c.Delta > 0 {
			upMoves++
		} else if c.Delta < 0 {
			downMoves++
		}
		memberMoves[c.MemberID]++
	}
	var mostActiveMember string
	maxMoves := 0
	for m, count := range memberMoves {
		if count > maxMoves {
			maxMoves = count
			mostActiveMember = m
		}
	}
	return map[string]interface{}{
		"board_id":           boardID,
		"total_moves":        totalMoves,
		"up_moves":           upMoves,
		"down_moves":         downMoves,
		"most_active_member": mostActiveMember,
		"most_active_count":  maxMoves,
	}, nil
}
