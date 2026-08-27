package service

import (
	"sort"
	"time"

	"leaderboard/internal/model"
)

func (s *Service) ListChangeLogs(filter model.ChangeLogFilter, page, size int) ([]*model.ChangeLog, int, error) {
	all := s.store.ListChangeLogs()
	matched := make([]*model.ChangeLog, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].At.After(matched[j].At)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ChangeLog{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetChangeLog(id string) (*model.ChangeLog, error) {
	return s.store.GetChangeLog(id)
}

func (s *Service) DeleteChangeLog(id string) error {
	return s.store.DeleteChangeLog(id)
}

type BoardReport struct {
	BoardID    string  `json:"board_id"`
	SeasonID   string  `json:"season_id"`
	MemberCount int    `json:"member_count"`
	TopScore   int64   `json:"top_score"`
	AvgScore   float64 `json:"avg_score"`
	MedianScore float64 `json:"median_score"`
	P90Score    float64 `json:"p90_score"`
}

func (s *Service) BoardReport(boardID, seasonID string) (*BoardReport, error) {
	all := s.store.ListRankEntries()
	var scores []int64
	for _, r := range all {
		if r.BoardID == boardID && r.SeasonID == seasonID {
			scores = append(scores, r.Score)
		}
	}
	if len(scores) == 0 {
		return &BoardReport{BoardID: boardID, SeasonID: seasonID}, nil
	}
	sort.Slice(scores, func(i, j int) bool { return scores[i] < scores[j] })
	var sum int64
	for _, v := range scores {
		sum += v
	}
	avg := float64(sum) / float64(len(scores))
	median := percentile(scores, 0.5)
	p90 := percentile(scores, 0.9)
	return &BoardReport{
		BoardID:     boardID,
		SeasonID:    seasonID,
		MemberCount: len(scores),
		TopScore:    scores[len(scores)-1],
		AvgScore:    avg,
		MedianScore: median,
		P90Score:    p90,
	}, nil
}

func percentile(sorted []int64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := float64(len(sorted)-1) * p
	lower := int(idx)
	upper := lower + 1
	if upper >= len(sorted) {
		return float64(sorted[lower])
	}
	frac := idx - float64(lower)
	return float64(sorted[lower])*(1-frac) + float64(sorted[upper])*frac
}

func (s *Service) ExportBoardData(boardID string) (map[string]interface{}, error) {
	board, err := s.store.GetBoard(boardID)
	if err != nil {
		return nil, err
	}
	seasons := s.store.ListSeasons()
	var boardSeasons []*model.Season
	for _, se := range seasons {
		if se.BoardID == boardID {
			boardSeasons = append(boardSeasons, se)
		}
	}
	entries := s.store.ListRankEntries()
	var boardEntries []*model.RankEntry
	for _, e := range entries {
		if e.BoardID == boardID {
			boardEntries = append(boardEntries, e)
		}
	}
	logs := s.store.ListChangeLogs()
	var boardLogs []*model.ChangeLog
	for _, c := range logs {
		if c.BoardID == boardID {
			boardLogs = append(boardLogs, c)
		}
	}
	return map[string]interface{}{
		"board":       board,
		"seasons":     boardSeasons,
		"entries":     boardEntries,
		"change_logs": boardLogs,
		"exported_at": time.Now(),
	}, nil
}
