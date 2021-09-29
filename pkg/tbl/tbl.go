package tbl

import (
	"context"
	"time"
)

// Tabler defines methods for adding event data to
// Google Sheets.
type Tabler interface {
	AppendRow(ctx context.Context, row Row) error
}

// Row represents a new row to append to the target
// Google Sheet.
type Row struct {
	ID            string     `json:"id"`
	ItemID        int        `json:"item_id"`
	Event         string     `json:"event"`
	UserID        int        `json:"user_id"`
	UserEmail     string     `json:"email"`
	ProjectID     int        `json:"project_id"`
	ProjectName   string     `json:"project_name"`
	Content       string     `json:"content"`
	Description   string     `json:"description"`
	Notes         []string   `json:"notes"`
	Priority      int        `json:"priority"`
	ParentID      *int       `json:"parent_id,omitempty"`
	SectionID     *int       `json:"section_id,omitempty"`
	SectionName   *string    `json:"section_name,omitempty"`
	LabelIDs      []int      `json:"label_ids"`
	LabelNames    []string   `json:"label_names"`
	Checked       bool       `json:"checked"`
	DateAdded     time.Time  `json:"date_added"`
	DateCompleted *time.Time `json:"date_completed,omitempty"`
}
