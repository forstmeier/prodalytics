package tbl

import (
	"context"
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
	if err := c.helper.appendRow(ctx, row); err != nil {
		return &AppendRowError{
			err: err,
		}
	}

	return nil
}
