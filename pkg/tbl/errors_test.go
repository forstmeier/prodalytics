package tbl

import (
	"errors"
	"testing"
)

func TestErrorAppendRow(t *testing.T) {
	err := &ErrorAppendRow{
		err: errors.New("mock append row error"),
	}

	recieved := err.Error()
	expected := "[tbl] [append row]: mock append row error"

	if recieved != expected {
		t.Errorf("incorrect error message, received: %s, expected: %s", recieved, expected)
	}
}
