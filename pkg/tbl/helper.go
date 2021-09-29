package tbl

import (
	"context"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"
)

var _ helper = &help{}

type helper interface {
	appendRow(ctx context.Context, row Row) error
}

type help struct {
	sheetID      string
	sheetsClient sheetsClient
}

type sheetsClient interface {
	Append(spreadsheetID string, startRange string, valueRange *sheets.ValueRange) appendCaller
}

// appendCaller is included for test mocking
type appendCaller interface {
	Do(opts ...googleapi.CallOption) (*sheets.AppendValuesResponse, error)
}

type appendCall struct {
	err error
}

func (a *appendCall) Do(opts ...googleapi.CallOption) (*sheets.AppendValuesResponse, error) {
	return nil, a.err
}

func (h *help) appendRow(ctx context.Context, row Row) error {
	valueRange := &sheets.ValueRange{
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

	_, err := h.sheetsClient.Append(h.sheetID, "A1", valueRange).Do()
	if err != nil {
		return err
	}

	return nil
}
