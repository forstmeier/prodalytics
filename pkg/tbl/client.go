package tbl

import (
	"context"
	"strconv"
	"strings"

	"google.golang.org/api/sheets/v4"
)

var _ Tabler = &Client{}

// Client implements the tbl.Tabler methods using
// Google Sheets.
type Client struct {
	helper helper
}

// New generates a Client pointer instance.
func New(sheetID string, sheetsClient sheetsClient) *Client {
	return &Client{
		helper: &help{
			sheetID:      sheetID,
			sheetsClient: sheetsClient,
		},
	}
}

// AppendRow implements the tbl.Tabler.AppendRow
// interface method.
func (c *Client) AppendRow(ctx context.Context, row Row) error {
	labelIDs := ""
	for _, labelID := range row.LabelIDs {
		labelIDs += strconv.Itoa(labelID)
	}

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
				strings.Join(row.Notes, "\n\n"),
				row.Priority,
				row.ParentID,
				row.SectionID,
				row.SectionName,
				labelIDs,
				strings.Join(row.LabelNames, ","),
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
