package service

import (
	"sort"
	"time"

	"leaderboard/internal/model"
	"leaderboard/pkg/idgen"
)

type LeaderboardOverview struct {
	Board          *model.Board     `json:"board"`
	ActiveSeason   *model.Season    `json:"active_season,omitempty"`
	TopEntries     []*model.RankEntry `json:"top_entries"`
	TotalMembers   int              `json:"total_members"`
	TotalEvents    int              `json:"total_events"`
}

func (s *Service) GetLeaderboardOverview(boardID string) (*LeaderboardOverview, error) {
	board, err := s.store.GetBoard(boardID)
	if err != nil {
		return nil, err
	}
	seasons := s.store.ListSeasons()
	var activeSeason *model.Season
	for _, se := range seasons {
		if se.BoardID == boardID && se.Status == model.SeasonStatusActive {
			activeSeason = se
			break
		}
	}
	seasonID := ""
	if activeSeason != nil {
		seasonID = activeSeason.ID
	}
	topEntries, _ := s.GetTopN(boardID, seasonID, 10)

	allMembers := s.store.ListMembers()
	allEvents := s.store.ListScoreEvents()
	var memberCount, eventCount int
	for _, e := range allEvents {
		if e.BoardID == boardID {
			eventCount++
		}
	}
	memberIDs := make(map[string]bool)
	entries := s.store.ListRankEntries()
	for _, r := range entries {
		if r.BoardID == boardID && r.SeasonID == seasonID {
			memberIDs[r.MemberID] = true
		}
	}
	memberCount = len(memberIDs)
	if memberCount == 0 {
		memberCount = len(allMembers)
	}

	return &LeaderboardOverview{
		Board:        board,
		ActiveSeason: activeSeason,
		TopEntries:   topEntries,
		TotalMembers: memberCount,
		TotalEvents:  eventCount,
	}, nil
}

type MemberLeaderboardProfile struct {
	Member       *model.Member      `json:"member"`
	Boards       []BoardProfile     `json:"boards"`
	TotalScore   int64              `json:"total_score"`
	BestRank     int                `json:"best_rank"`
	RankHistory  []*model.RankEntry `json:"rank_history"`
}

type BoardProfile struct {
	BoardID   string `json:"board_id"`
	BoardName string `json:"board_name"`
	Score     int64  `json:"score"`
	Rank      int    `json:"rank"`
	SeasonID  string `json:"season_id"`
}

func (s *Service) GetMemberLeaderboardProfile(memberID string) (*MemberLeaderboardProfile, error) {
	member, err := s.store.GetMember(memberID)
	if err != nil {
		return nil, err
	}
	entries := s.store.ListRankEntries()
	var boards []BoardProfile
	var totalScore int64
	bestRank := 0
	var history []*model.RankEntry
	for _, r := range entries {
		if r.MemberID == memberID {
			history = append(history, r)
			totalScore += r.Score
			if bestRank == 0 || r.Rank < bestRank {
				bestRank = r.Rank
			}
			board, _ := s.store.GetBoard(r.BoardID)
			name := ""
			if board != nil {
				name = board.Name
			}
			boards = append(boards, BoardProfile{
				BoardID:   r.BoardID,
				BoardName: name,
				Score:     r.Score,
				Rank:      r.Rank,
				SeasonID:  r.SeasonID,
			})
		}
	}
	sort.Slice(history, func(i, j int) bool {
		return history[i].UpdatedAt.After(history[j].UpdatedAt)
	})
	return &MemberLeaderboardProfile{
		Member:      member,
		Boards:      boards,
		TotalScore:  totalScore,
		BestRank:    bestRank,
		RankHistory: history,
	}, nil
}

type SeasonComparison struct {
	SeasonID    string          `json:"season_id"`
	SeasonName  string          `json:"season_name"`
	MemberCount int             `json:"member_count"`
	TopScore    int64           `json:"top_score"`
	AvgScore    float64         `json:"avg_score"`
	PrevSeasonDelta float64     `json:"prev_season_delta"`
}

func (s *Service) CompareSeasons(boardID string) ([]SeasonComparison, error) {
	seasons := s.store.ListSeasons()
	var boardSeasons []*model.Season
	for _, se := range seasons {
		if se.BoardID == boardID {
			boardSeasons = append(boardSeasons, se)
		}
	}
	sort.Slice(boardSeasons, func(i, j int) bool {
		return boardSeasons[i].StartAt.Before(boardSeasons[j].StartAt)
	})

	entries := s.store.ListRankEntries()
	var result []SeasonComparison
	var prevAvg float64
	for _, se := range boardSeasons {
		var scores []int64
		var count int
		for _, r := range entries {
			if r.BoardID == boardID && r.SeasonID == se.ID {
				scores = append(scores, r.Score)
				count++
			}
		}
		if len(scores) == 0 {
			result = append(result, SeasonComparison{
				SeasonID:    se.ID,
				SeasonName:  se.Name,
				MemberCount: 0,
			})
			continue
		}
		sort.Slice(scores, func(i, j int) bool { return scores[i] < scores[j] })
		var sum int64
		for _, v := range scores {
			sum += v
		}
		avg := float64(sum) / float64(len(scores))
		delta := 0.0
		if prevAvg > 0 {
			delta = avg - prevAvg
		}
		result = append(result, SeasonComparison{
			SeasonID:        se.ID,
			SeasonName:      se.Name,
			MemberCount:     count,
			TopScore:        scores[len(scores)-1],
			AvgScore:        avg,
			PrevSeasonDelta: delta,
		})
		prevAvg = avg
	}
	return result, nil
}

type ActivityTrend struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
	Score int64  `json:"score"`
}

func (s *Service) GetBoardActivityTrend(boardID string, days int) ([]ActivityTrend, error) {
	if days < 1 {
		days = 7
	}
	if days > 90 {
		days = 90
	}
	allEvents := s.store.ListScoreEvents()
	cutoff := time.Now().AddDate(0, 0, -days)
	dayMap := make(map[string]*ActivityTrend)
	for _, e := range allEvents {
		if e.BoardID != boardID || e.At.Before(cutoff) {
			continue
		}
		dateKey := e.At.Format("2006-01-02")
		if _, ok := dayMap[dateKey]; !ok {
			dayMap[dateKey] = &ActivityTrend{Date: dateKey}
		}
		dayMap[dateKey].Count++
		dayMap[dateKey].Score += e.Value
	}
	var result []ActivityTrend
	for d := 0; d < days; d++ {
		dateKey := time.Now().AddDate(0, 0, -d).Format("2006-01-02")
		if v, ok := dayMap[dateKey]; ok {
			result = append(result, *v)
		} else {
			result = append(result, ActivityTrend{Date: dateKey, Count: 0, Score: 0})
		}
	}
	return result, nil
}

func (s *Service) BulkCreateMembers(inputs []model.Member) ([]*model.Member, error) {
	var result []*model.Member
	for _, input := range inputs {
		m, err := s.CreateMember(input)
		if err != nil {
			return result, err
		}
		result = append(result, m)
	}
	return result, nil
}

func (s *Service) BulkCreateScoreEvents(inputs []model.ScoreEvent) ([]*model.ScoreEvent, error) {
	var result []*model.ScoreEvent
	for _, input := range inputs {
		e, err := s.CreateScoreEvent(input)
		if err != nil {
			return result, err
		}
		result = append(result, e)
	}
	return result, nil
}

func (s *Service) RecalculateAllRanks(boardID, seasonID string) error {
	board, err := s.store.GetBoard(boardID)
	if err != nil {
		return err
	}
	allEvents := s.store.ListScoreEvents()
	memberScores := make(map[string]int64)
	for _, e := range allEvents {
		if e.BoardID == boardID {
			if seasonID == "" || e.SeasonID == seasonID {
				memberScores[e.MemberID] += e.Value
			}
		}
	}
	var entries []*model.RankEntry
	for memberID, score := range memberScores {
		entry, err := s.store.GetRankEntryByBoardSeasonMember(boardID, seasonID, memberID)
		now := time.Now()
		if err != nil {
			entry = &model.RankEntry{
				ID:        idgen.Hex(),
				BoardID:   boardID,
				SeasonID:  seasonID,
				MemberID:  memberID,
				Score:     score,
				UpdatedAt: now,
			}
			_ = s.store.CreateRankEntry(entry)
		} else {
			entry.Score = score
			entry.UpdatedAt = now
			_ = s.store.UpdateRankEntry(entry)
		}
		entries = append(entries, entry)
	}
	asc := board.SortOrder == model.BoardSortAsc
	sort.Slice(entries, func(i, j int) bool {
		if asc {
			return entries[i].Score < entries[j].Score
		}
		return entries[i].Score > entries[j].Score
	})
	for i, r := range entries {
		r.PrevRank = r.Rank
		r.Rank = i + 1
		_ = s.store.UpdateRankEntry(r)
		if r.Rank != r.PrevRank && r.PrevRank != 0 {
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

type SystemStats struct {
	TotalBoards        int `json:"total_boards"`
	TotalMembers       int `json:"total_members"`
	TotalSeasons       int `json:"total_seasons"`
	TotalScoreEvents   int `json:"total_score_events"`
	TotalRankEntries   int `json:"total_rank_entries"`
	TotalSnapshots     int `json:"total_snapshots"`
	TotalChangeLogs    int `json:"total_change_logs"`
	ActiveBoards       int `json:"active_boards"`
	ActiveSeasons      int `json:"active_seasons"`
}

func (s *Service) GetSystemStats() *SystemStats {
	boards := s.store.ListBoards()
	members := s.store.ListMembers()
	seasons := s.store.ListSeasons()
	events := s.store.ListScoreEvents()
	entries := s.store.ListRankEntries()
	snapshots := s.store.ListRankSnapshots()
	logs := s.store.ListChangeLogs()

	var activeBoards, activeSeasons int
	for _, b := range boards {
		if b.Status == model.BoardStatusActive {
			activeBoards++
		}
	}
	for _, se := range seasons {
		if se.Status == model.SeasonStatusActive {
			activeSeasons++
		}
	}

	return &SystemStats{
		TotalBoards:      len(boards),
		TotalMembers:     len(members),
		TotalSeasons:     len(seasons),
		TotalScoreEvents: len(events),
		TotalRankEntries: len(entries),
		TotalSnapshots:   len(snapshots),
		TotalChangeLogs:  len(logs),
		ActiveBoards:     activeBoards,
		ActiveSeasons:    activeSeasons,
	}
}
