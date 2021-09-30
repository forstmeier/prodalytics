package tbl

import (
	"context"
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	sheetsClient := &mockSheetsClient{}

	client := New(sheetsClient)
	if client == nil {
		t.Error("error creating tbl client")
	}
}

func TestAppendRow(t *testing.T) {
	tests := []struct {
		description        string
		mockAppendRowError error
		error              error
	}{
		{
			description:        "error append row",
			mockAppendRowError: errors.New("mock append row error"),
			error:              &AppendRowError{},
		},
		{
			description:        "successful append row invocation",
			mockAppendRowError: nil,
			error:              nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			client := &Client{
				helper: &help{
					sheetsClient: &mockSheetsClient{
						mockAppendRowError: test.mockAppendRowError,
					},
				},
			}

			err := client.AppendRow(context.Background(), Row{})

			if err != nil {
				switch e := test.error.(type) {
				case *AppendRowError:
					if !errors.As(err, &e) {
						t.Errorf("incorrect error, received: %v, expected: %v", err, e)
					}
				default:
					t.Fatalf("unexpected error type: %v", err)
				}
			}
		})
	}
}
