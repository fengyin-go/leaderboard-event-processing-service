package model

import (
	"strings"
	"time"
)

const (
	BoardMetricScore  = "score"
	BoardMetricStreak = "streak"
	BoardMetricCount  = "count"
)

const (
	BoardSortAsc  = "asc"
	BoardSortDesc = "desc"
)

const (
	BoardStatusActive   = "active"
	BoardStatusArchived = "archived"
)

type Board struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Metric      string    `json:"metric"`
	SortOrder   string    `json:"sort_order"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (b *Board) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	b.Description = strings.TrimSpace(b.Description)
	b.Metric = strings.TrimSpace(b.Metric)
	b.SortOrder = strings.TrimSpace(b.SortOrder)
	b.Status = strings.TrimSpace(b.Status)

	if b.Name == "" {
		return NewValidationError("name", "榜单名称不能为空")
	}
	if b.Metric == "" {
		b.Metric = BoardMetricScore
	}
	if b.Metric != BoardMetricScore && b.Metric != BoardMetricStreak && b.Metric != BoardMetricCount {
		return NewValidationError("metric", "指标类型不合法")
	}
	if b.SortOrder == "" {
		b.SortOrder = BoardSortDesc
	}
	if b.SortOrder != BoardSortAsc && b.SortOrder != BoardSortDesc {
		return NewValidationError("sort_order", "排序方向不合法")
	}
	if b.Status == "" {
		b.Status = BoardStatusActive
	}
	if b.Status != BoardStatusActive && b.Status != BoardStatusArchived {
		return NewValidationError("status", "榜单状态不合法")
	}
	return nil
}

type BoardFilter struct {
	Name      string
	Metric    string
	Status    string
	SortOrder string
	Keyword   string
}

func (f BoardFilter) Match(b *Board) bool {
	if f.Name != "" && b.Name != f.Name {
		return false
	}
	if f.Metric != "" && b.Metric != f.Metric {
		return false
	}
	if f.Status != "" && b.Status != f.Status {
		return false
	}
	if f.SortOrder != "" && b.SortOrder != f.SortOrder {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(b.Name), k) && !strings.Contains(strings.ToLower(b.Description), k) {
			return false
		}
	}
	return true
}
