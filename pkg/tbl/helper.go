package tbl

import (
	"context"

	"google.golang.org/api/sheets/v4"
)

var _ helper = &help{}

type helper interface {
	appendRow(ctx context.Context, values *sheets.ValueRange) error
}

type help struct {
	sheetID      string
	sheetsClient sheetsClient
}

type sheetsClient interface {
	Append(spreadsheetID string, startRange string, valueRange *sheets.ValueRange) *sheets.SpreadsheetsValuesAppendCall
}

func (h *help) appendRow(ctx context.Context, values *sheets.ValueRange) error {
	_, err := h.sheetsClient.Append(h.sheetID, "A1", values).ValueInputOption("RAW").Do()
	if err != nil {
		return err
	}

	return nil
}
