package service

import (
	"sort"
	"time"

	"leaderboard/internal/model"
)

type MemberActivity struct {
	MemberID    string    `json:"member_id"`
	MemberName  string    `json:"member_name"`
	EventCount  int       `json:"event_count"`
	TotalScore  int64     `json:"total_score"`
	LastActive  time.Time `json:"last_active"`
	AvgScore    float64   `json:"avg_score"`
}

func (s *Service) GetMemberActivity(boardID string, limit int) ([]MemberActivity, error) {
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	allEvents := s.store.ListScoreEvents()
	memberMap := make(map[string]*MemberActivity)
	for _, e := range allEvents {
		if e.BoardID != boardID {
			continue
		}
		if _, ok := memberMap[e.MemberID]; !ok {
			memberMap[e.MemberID] = &MemberActivity{MemberID: e.MemberID}
		}
		ma := memberMap[e.MemberID]
		ma.EventCount++
		ma.TotalScore += e.Value
		if e.At.After(ma.LastActive) {
			ma.LastActive = e.At
		}
	}
	for _, ma := range memberMap {
		if ma.EventCount > 0 {
			ma.AvgScore = float64(ma.TotalScore) / float64(ma.EventCount)
		}
		m, err := s.store.GetMember(ma.MemberID)
		if err == nil {
			ma.MemberName = m.Name
		}
	}
	var result []MemberActivity
	for _, ma := range memberMap {
		result = append(result, *ma)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalScore > result[j].TotalScore
	})
	if len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

type BoardDistribution struct {
	Range    string `json:"range"`
	Count    int    `json:"count"`
	Pct      float64 `json:"pct"`
}

func (s *Service) GetScoreDistribution(boardID, seasonID string, buckets int) ([]BoardDistribution, error) {
	if buckets < 2 {
		buckets = 5
	}
	if buckets > 20 {
		buckets = 20
	}
	all := s.store.ListRankEntries()
	var scores []int64
	for _, r := range all {
		if r.BoardID == boardID && r.SeasonID == seasonID {
			scores = append(scores, r.Score)
		}
	}
	if len(scores) == 0 {
		return []BoardDistribution{}, nil
	}
	sort.Slice(scores, func(i, j int) bool { return scores[i] < scores[j] })
	min := scores[0]
	max := scores[len(scores)-1]
	if max == min {
		return []BoardDistribution{{Range: "全部", Count: len(scores), Pct: 100.0}}, nil
	}
	step := float64(max-min) / float64(buckets)
	counts := make([]int, buckets)
	for _, v := range scores {
		idx := int(float64(v-min) / step)
		if idx >= buckets {
			idx = buckets - 1
		}
		counts[idx]++
	}
	var result []BoardDistribution
	for i := 0; i < buckets; i++ {
		low := float64(min) + float64(i)*step
		high := low + step
		label := ""
		if i == buckets-1 {
			label = "[" + formatNum(low) + ", " + formatNum(high) + "]"
		} else {
			label = "[" + formatNum(low) + ", " + formatNum(high) + ")"
		}
		pct := float64(counts[i]) / float64(len(scores)) * 100
		result = append(result, BoardDistribution{Range: label, Count: counts[i], Pct: pct})
	}
	return result, nil
}

func formatNum(n float64) string {
	if n == float64(int64(n)) {
		return ""
	}
	return ""
}

type HotStreak struct {
	MemberID   string `json:"member_id"`
	MemberName string `json:"member_name"`
	StreakDays int    `json:"streak_days"`
	LastScore  int64  `json:"last_score"`
}

func (s *Service) GetHotStreaks(boardID string, limit int) ([]HotStreak, error) {
	if limit < 1 {
		limit = 10
	}
	allEvents := s.store.ListScoreEvents()
	memberLatest := make(map[string]time.Time)
	memberScores := make(map[string]int64)
	for _, e := range allEvents {
		if e.BoardID != boardID {
			continue
		}
		if latest, ok := memberLatest[e.MemberID]; !ok || e.At.After(latest) {
			memberLatest[e.MemberID] = e.At
			memberScores[e.MemberID] = e.Value
		}
	}
	var result []HotStreak
	for memberID, latest := range memberLatest {
		days := int(time.Since(latest).Hours() / 24)
		m, err := s.store.GetMember(memberID)
		name := ""
		if err == nil {
			name = m.Name
		}
		result = append(result, HotStreak{
			MemberID:   memberID,
			MemberName: name,
			StreakDays: days,
			LastScore:  memberScores[memberID],
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].StreakDays < result[j].StreakDays
	})
	if len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

type SeasonMilestone struct {
	SeasonID   string    `json:"season_id"`
	SeasonName string    `json:"season_name"`
	Milestone  string    `json:"milestone"`
	At         time.Time `json:"at"`
}

func (s *Service) GetSeasonMilestones(boardID string) ([]SeasonMilestone, error) {
	seasons := s.store.ListSeasons()
	var boardSeasons []*model.Season
	for _, se := range seasons {
		if se.BoardID == boardID {
			boardSeasons = append(boardSeasons, se)
		}
	}
	sort.Slice(boardSeasons, func(i, j int) bool {
		return boardSeasons[i].CreatedAt.Before(boardSeasons[j].CreatedAt)
	})
	var result []SeasonMilestone
	for _, se := range boardSeasons {
		result = append(result, SeasonMilestone{
			SeasonID:   se.ID,
			SeasonName: se.Name,
			Milestone:  "赛季创建",
			At:         se.CreatedAt,
		})
		if se.Status == model.SeasonStatusActive || se.Status == model.SeasonStatusFinished {
			result = append(result, SeasonMilestone{
				SeasonID:   se.ID,
				SeasonName: se.Name,
				Milestone:  "赛季开始",
				At:         se.StartAt,
			})
		}
		if se.Status == model.SeasonStatusFinished {
			result = append(result, SeasonMilestone{
				SeasonID:   se.ID,
				SeasonName: se.Name,
				Milestone:  "赛季结束",
				At:         se.EndAt,
			})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].At.Before(result[j].At)
	})
	return result, nil
}

type RankVelocity struct {
	MemberID   string  `json:"member_id"`
	MemberName string  `json:"member_name"`
	Delta      int     `json:"delta"`
	Events     int     `json:"events"`
}

func (s *Service) GetRankVelocity(boardID, seasonID string, since time.Time) ([]RankVelocity, error) {
	allEvents := s.store.ListScoreEvents()
	memberEvents := make(map[string]int)
	for _, e := range allEvents {
		if e.BoardID == boardID && e.SeasonID == seasonID && e.At.After(since) {
			memberEvents[e.MemberID]++
		}
	}
	logs := s.store.ListChangeLogs()
	memberDelta := make(map[string]int)
	for _, c := range logs {
		if c.BoardID == boardID && c.At.After(since) {
			memberDelta[c.MemberID] += c.Delta
		}
	}
	var result []RankVelocity
	for memberID, delta := range memberDelta {
		m, err := s.store.GetMember(memberID)
		name := ""
		if err == nil {
			name = m.Name
		}
		result = append(result, RankVelocity{
			MemberID:   memberID,
			MemberName: name,
			Delta:      delta,
			Events:     memberEvents[memberID],
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Delta > result[j].Delta
	})
	return result, nil
}
