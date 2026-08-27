package service

import (
	"sort"

	"leaderboard/internal/model"
)

func (s *Service) SearchMembersByTag(tag string) (*model.Member, error) {
	return s.store.GetMemberByTag(tag)
}

func (s *Service) ListMembersByIDs(ids []string) []*model.Member {
	all := s.store.ListMembers()
	idMap := make(map[string]bool)
	for _, id := range ids {
		idMap[id] = true
	}
	var result []*model.Member
	for _, m := range all {
		if idMap[m.ID] {
			result = append(result, m)
		}
	}
	return result
}

func (s *Service) GetMemberScoreEvents(memberID string, limit int) ([]*model.ScoreEvent, error) {
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	all := s.store.ListScoreEvents()
	var result []*model.ScoreEvent
	for _, e := range all {
		if e.MemberID == memberID {
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

func (s *Service) GetMemberTotalScore(memberID string) int64 {
	all := s.store.ListScoreEvents()
	var total int64
	for _, e := range all {
		if e.MemberID == memberID {
			total += e.Value
		}
	}
	return total
}

func (s *Service) GetMemberRankHistory(memberID, boardID string) ([]*model.RankEntry, error) {
	all := s.store.ListRankEntries()
	var result []*model.RankEntry
	for _, r := range all {
		if r.MemberID == memberID && r.BoardID == boardID {
			result = append(result, r)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].UpdatedAt.After(result[j].UpdatedAt)
	})
	return result, nil
}

func (s *Service) MergeMembers(fromID, toID string) error {
	from, err := s.store.GetMember(fromID)
	if err != nil {
		return err
	}
	to, err := s.store.GetMember(toID)
	if err != nil {
		return err
	}
	_ = from
	_ = to

	events := s.store.ListScoreEvents()
	for _, e := range events {
		if e.MemberID == fromID {
			e.MemberID = toID
		}
	}
	entries := s.store.ListRankEntries()
	for _, r := range entries {
		if r.MemberID == fromID {
			r.MemberID = toID
		}
	}
	logs := s.store.ListChangeLogs()
	for _, c := range logs {
		if c.MemberID == fromID {
			c.MemberID = toID
		}
	}
	return s.store.DeleteMember(fromID)
}
