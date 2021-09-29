package tbl

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"google.golang.org/api/sheets/v4"
)

type mockSheetsClient struct {
	rowID              string
	mockAppendRowError error
}

func (m *mockSheetsClient) Append(spreadsheetID string, startRange string, valueRange *sheets.ValueRange) appendCaller {
	m.rowID = valueRange.Values[0][0].(string)

	return &appendCall{
		err: m.mockAppendRowError,
	}
}

func Test_appendRow(t *testing.T) {
	tests := []struct {
		description        string
		mockAppendRowError error
		error              string
	}{
		{
			description:        "error append",
			mockAppendRowError: errors.New("mock append error"),
			error:              "mock append error",
		},
		{
			description:        "successful append row invocation",
			mockAppendRowError: nil,
			error:              "",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			h := &help{
				sheetID: "sheetID",
				sheetsClient: &mockSheetsClient{
					mockAppendRowError: test.mockAppendRowError,
				},
			}

			id := uuid.NewString()

			err := h.appendRow(context.Background(), Row{ID: id})

			if err != nil {
				if err.Error() != test.error {
					t.Errorf("incorrect error, received: %s, expected: %s", err.Error(), test.error)
				}
			} else {
				received := h.sheetsClient.(*mockSheetsClient).rowID
				expected := id
				if received != expected {
					t.Errorf("incorrect row id, received: %s, expected: %s", received, expected)
				}
			}
		})
	}
}
