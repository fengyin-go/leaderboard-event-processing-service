package service

import (
	"sort"

	"leaderboard/internal/model"
)

func (s *Service) GetRankEntry(id string) (*model.RankEntry, error) {
	return s.store.GetRankEntry(id)
}

func (s *Service) GetMemberRank(boardID, seasonID, memberID string) (*model.RankEntry, error) {
	return s.store.GetRankEntryByBoardSeasonMember(boardID, seasonID, memberID)
}

func (s *Service) ListRankEntries(filter model.RankEntryFilter, page, size int) ([]*model.RankEntry, int, error) {
	all := s.store.ListRankEntries()
	matched := make([]*model.RankEntry, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Rank < matched[j].Rank
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RankEntry{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetTopN(boardID, seasonID string, n int) ([]*model.RankEntry, error) {
	all := s.store.ListRankEntries()
	var matched []*model.RankEntry
	for _, r := range all {
		if r.BoardID == boardID && r.SeasonID == seasonID {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Rank < matched[j].Rank
	})
	if n > len(matched) {
		n = len(matched)
	}
	return matched[:n], nil
}

func (s *Service) GetRankRange(boardID, seasonID string, startRank, endRank int) ([]*model.RankEntry, error) {
	all := s.store.ListRankEntries()
	var matched []*model.RankEntry
	for _, r := range all {
		if r.BoardID == boardID && r.SeasonID == seasonID && r.Rank >= startRank && r.Rank <= endRank {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Rank < matched[j].Rank
	})
	return matched, nil
}
