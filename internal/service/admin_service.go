package service

import (
	"sort"
	"time"

	"leaderboard/internal/model"
	"leaderboard/pkg/idgen"
)

func (s *Service) AdminListAllBoards(page, size int) ([]*model.Board, int, error) {
	return s.ListBoards(model.BoardFilter{}, page, size)
}

func (s *Service) AdminPurgeBoardData(boardID string) error {
	if _, err := s.store.GetBoard(boardID); err != nil {
		return err
	}
	events := s.store.ListScoreEvents()
	for _, e := range events {
		if e.BoardID == boardID {
			_ = s.store.DeleteScoreEvent(e.ID)
		}
	}
	entries := s.store.ListRankEntries()
	for _, r := range entries {
		if r.BoardID == boardID {
			_ = s.store.DeleteRankEntry(r.ID)
		}
	}
	logs := s.store.ListChangeLogs()
	for _, c := range logs {
		if c.BoardID == boardID {
			_ = s.store.DeleteChangeLog(c.ID)
		}
	}
	snapshots := s.store.ListRankSnapshots()
	for _, rs := range snapshots {
		if rs.BoardID == boardID {
			_ = s.store.DeleteRankSnapshot(rs.ID)
		}
	}
	seasons := s.store.ListSeasons()
	for _, se := range seasons {
		if se.BoardID == boardID {
			_ = s.store.DeleteSeason(se.ID)
		}
	}
	return s.store.DeleteBoard(boardID)
}

func (s *Service) AdminGlobalReport() (map[string]interface{}, error) {
	boards := s.store.ListBoards()
	members := s.store.ListMembers()
	seasons := s.store.ListSeasons()
	entries := s.store.ListRankEntries()
	events := s.store.ListScoreEvents()
	var totalScore int64
	for _, e := range events {
		totalScore += e.Value
	}
	var avgScore float64
	if len(entries) > 0 {
		var sum int64
		for _, r := range entries {
			sum += r.Score
		}
		avgScore = float64(sum) / float64(len(entries))
	}
	return map[string]interface{}{
		"total_boards":    len(boards),
		"total_members":   len(members),
		"total_seasons":   len(seasons),
		"total_entries":   len(entries),
		"total_events":    len(events),
		"total_score":     totalScore,
		"avg_entry_score": avgScore,
		"generated_at":    time.Now(),
	}, nil
}

func (s *Service) AdminListInactiveMembers(since time.Time) ([]*model.Member, error) {
	all := s.store.ListMembers()
	var result []*model.Member
	for _, m := range all {
		if m.UpdatedAt.Before(since) && m.CreatedAt.Before(since) {
			result = append(result, m)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result, nil
}

func (s *Service) AdminRebuildAllChangeLogs(boardID string) error {
	logs := s.store.ListChangeLogs()
	for _, c := range logs {
		if c.BoardID == boardID {
			_ = s.store.DeleteChangeLog(c.ID)
		}
	}
	entries := s.store.ListRankEntries()
	for _, r := range entries {
		if r.BoardID == boardID && r.PrevRank != 0 && r.PrevRank != r.Rank {
			cl := &model.ChangeLog{
				ID:       idgen.Hex(),
				BoardID:  r.BoardID,
				MemberID: r.MemberID,
				OldRank:  r.PrevRank,
				NewRank:  r.Rank,
				Delta:    r.PrevRank - r.Rank,
				At:       time.Now(),
			}
			_ = s.store.CreateChangeLog(cl)
		}
	}
	return nil
}
