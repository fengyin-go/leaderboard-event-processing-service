package service

import (
	"sort"
	"time"

	"leaderboard/internal/model"
	"leaderboard/internal/store"
	"leaderboard/pkg/idgen"
)

func (s *Service) CreateScoreEvent(input model.ScoreEvent) (*model.ScoreEvent, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetBoard(input.BoardID); err != nil {
		return nil, model.NewValidationError("board_id", "所属榜单不存在")
	}
	if _, err := s.store.GetMember(input.MemberID); err != nil {
		return nil, model.NewValidationError("member_id", "所属成员不存在")
	}
	if input.SeasonID != "" {
		if _, err := s.store.GetSeason(input.SeasonID); err != nil {
			return nil, model.NewValidationError("season_id", "所属赛季不存在")
		}
	}
	now := time.Now()
	e := &model.ScoreEvent{
		ID:       idgen.Hex(),
		BoardID:  input.BoardID,
		MemberID: input.MemberID,
		SeasonID: input.SeasonID,
		Value:    input.Value,
		Source:   input.Source,
		At:       now,
	}
	if err := s.store.CreateScoreEvent(e); err != nil {
		return nil, err
	}
	if err := s.recalcRankEntry(input.BoardID, input.SeasonID, input.MemberID); err != nil {
		s.log.Warnf("重算排名失败: %v", err)
	}
	return e, nil
}

func (s *Service) ListScoreEvents(filter model.ScoreEventFilter, page, size int) ([]*model.ScoreEvent, int, error) {
	all := s.store.ListScoreEvents()
	matched := make([]*model.ScoreEvent, 0, len(all))
	for _, e := range all {
		if filter.Match(e) {
			matched = append(matched, e)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].At.After(matched[j].At)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ScoreEvent{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetScoreEvent(id string) (*model.ScoreEvent, error) {
	return s.store.GetScoreEvent(id)
}

func (s *Service) DeleteScoreEvent(id string) error {
	return s.store.DeleteScoreEvent(id)
}

func (s *Service) recalcRankEntry(boardID, seasonID, memberID string) error {
	board, err := s.store.GetBoard(boardID)
	if err != nil {
		return err
	}
	allEvents := s.store.ListScoreEvents()
	var total int64
	for _, e := range allEvents {
		if e.BoardID == boardID && e.MemberID == memberID {
			if seasonID == "" || e.SeasonID == seasonID {
				total += e.Value
			}
		}
	}

	entry, err := s.store.GetRankEntryByBoardSeasonMember(boardID, seasonID, memberID)
	if err != nil && err != store.ErrNotFound {
		return err
	}
	now := time.Now()
	oldRank := 0
	if entry != nil {
		oldRank = entry.Rank
		entry.Score = total
		entry.UpdatedAt = now
	} else {
		entry = &model.RankEntry{
			ID:        idgen.Hex(),
			BoardID:   boardID,
			SeasonID:  seasonID,
			MemberID:  memberID,
			Score:     total,
			Rank:      0,
			PrevRank:  0,
			UpdatedAt: now,
		}
		if err := s.store.CreateRankEntry(entry); err != nil {
			return err
		}
	}

	allEntries := s.store.ListRankEntries()
	var relevant []*model.RankEntry
	for _, r := range allEntries {
		if r.BoardID == boardID && r.SeasonID == seasonID {
			relevant = append(relevant, r)
		}
	}

	asc := board.SortOrder == model.BoardSortAsc
	sort.Slice(relevant, func(i, j int) bool {
		if asc {
			return relevant[i].Score < relevant[j].Score
		}
		return relevant[i].Score > relevant[j].Score
	})

	for i, r := range relevant {
		r.PrevRank = r.Rank
		r.Rank = i + 1
		if err := s.store.UpdateRankEntry(r); err != nil {
			return err
		}
		if r.Rank != r.PrevRank && r.PrevRank != 0 {
			cl := &model.ChangeLog{
				ID:       idgen.Hex(),
				BoardID:  r.BoardID,
				MemberID: r.MemberID,
				OldRank:  r.PrevRank,
				NewRank:  r.Rank,
				Delta:    r.PrevRank - r.Rank,
				At:       now,
			}
			_ = s.store.CreateChangeLog(cl)
		}
	}

	if oldRank != 0 && entry.Rank != oldRank {
		cl := &model.ChangeLog{
			ID:       idgen.Hex(),
			BoardID:  entry.BoardID,
			MemberID: entry.MemberID,
			OldRank:  oldRank,
			NewRank:  entry.Rank,
			Delta:    oldRank - entry.Rank,
			At:       now,
		}
		_ = s.store.CreateChangeLog(cl)
	}

	return nil
}
