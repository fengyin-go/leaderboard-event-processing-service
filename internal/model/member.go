package model

import (
	"strings"
	"time"
)

type Member struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Tag       string    `json:"tag"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Member) Validate() error {
	m.Name = strings.TrimSpace(m.Name)
	m.Tag = strings.TrimSpace(m.Tag)

	if m.Name == "" {
		return NewValidationError("name", "成员名称不能为空")
	}
	if m.Tag == "" {
		return NewValidationError("tag", "成员标签不能为空")
	}
	return nil
}

type MemberFilter struct {
	Name    string
	Tag     string
	Keyword string
}

func (f MemberFilter) Match(m *Member) bool {
	if f.Name != "" && m.Name != f.Name {
		return false
	}
	if f.Tag != "" && m.Tag != f.Tag {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(m.Name), k) && !strings.Contains(strings.ToLower(m.Tag), k) {
			return false
		}
	}
	return true
}
