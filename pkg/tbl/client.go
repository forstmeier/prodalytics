package tbl

import (
	"context"

	"google.golang.org/api/sheets/v4"
)

var _ Tabler = &Client{}

// Client implements the tbl.Tabler methods using
// Google Sheets.
type Client struct {
	helper helper
}

// New generates a Client pointer instance.
func New(sheetsClient sheetsClient) *Client {
	return &Client{
		helper: &help{
			sheetsClient: sheetsClient,
		},
	}
}

// AppendRow implements the tbl.Tabler.AppendRow
// interface method.
func (c *Client) AppendRow(ctx context.Context, row Row) error {
	values := &sheets.ValueRange{
		Values: [][]interface{}{
			{
				row.ID,
				row.ItemID,
				row.Event,
				row.UserID,
				row.UserEmail,
				row.ProjectID,
				row.ProjectName,
				row.Content,
				row.Description,
				row.Notes,
				row.Priority,
				row.ParentID,
				row.SectionID,
				row.SectionName,
				row.LabelIDs,
				row.LabelNames,
				row.Checked,
				row.DateAdded,
				row.DateCompleted,
			},
		},
	}

	if err := c.helper.appendRow(ctx, values); err != nil {
		return &AppendRowError{
			err: err,
		}
	}

	return nil
}
