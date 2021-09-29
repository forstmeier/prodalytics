package evt

import (
	"context"

	"github.com/forstmeier/todalytics/pkg/tbl"
)

// Eventer defines methods for processing raw Todoist
// webhook events into table-formatted data.
type Eventer interface {
	Convert(ctx context.Context, data []byte) (*tbl.Row, error)
}

// Event represents a raw Todoist event.
type Event struct {
	EventName string    `json:"event_name"`
	EventData Item      `json:"event_data"`
	Initiator initiator `json:"initiator"`
}

// Item item represents an Item unmarshaled from the
// "event_data" field of a raw Todoist event.
//
// Note that only Item is being provided for the initial
// structure since these are the only events of interest.
type Item struct {
	ID            int     `json:"id"`
	UserID        int     `json:"user_id"`
	ProjectID     int     `json:"project_id"`
	Content       string  `json:"content"`
	Description   string  `json:"description"`
	Priority      int     `json:"priority"`
	ParentID      *int    `json:"parent_id,omitempty"`
	SectionID     *int    `json:"section_id,omitempty"`
	Labels        []int   `json:"labels"`
	Checked       int     `json:"checked"`
	DateAdded     string  `json:"date_added"`
	DateCompleted *string `json:"date_completed,omitempty"`
}

type initiator struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}
