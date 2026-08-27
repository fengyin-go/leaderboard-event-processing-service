package service

import (
	"encoding/json"
	"sort"
	"time"

	"leaderboard/internal/model"
	"leaderboard/pkg/idgen"
)

func (s *Service) CaptureRankSnapshot(boardID, seasonID string) (*model.RankSnapshot, error) {
	if _, err := s.store.GetBoard(boardID); err != nil {
		return nil, model.NewValidationError("board_id", "所属榜单不存在")
	}
	entries, err := s.GetTopN(boardID, seasonID, 1000)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(entries)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	rs := &model.RankSnapshot{
		ID:          idgen.Hex(),
		BoardID:     boardID,
		SeasonID:    seasonID,
		CapturedAt:  now,
		EntriesJSON: string(data),
		CreatedAt:   now,
	}
	if err := s.store.CreateRankSnapshot(rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func (s *Service) ListRankSnapshots(filter model.RankSnapshotFilter, page, size int) ([]*model.RankSnapshot, int, error) {
	all := s.store.ListRankSnapshots()
	matched := make([]*model.RankSnapshot, 0, len(all))
	for _, rs := range all {
		if filter.Match(rs) {
			matched = append(matched, rs)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CapturedAt.After(matched[j].CapturedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RankSnapshot{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetRankSnapshot(id string) (*model.RankSnapshot, error) {
	return s.store.GetRankSnapshot(id)
}

func (s *Service) DeleteRankSnapshot(id string) error {
	return s.store.DeleteRankSnapshot(id)
}
