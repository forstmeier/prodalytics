package tbl

import (
	"errors"
	"testing"
)

func TestAppendRowError(t *testing.T) {
	err := &AppendRowError{
		err: errors.New("mock append row error"),
	}

	recieved := err.Error()
	expected := "package tbl: mock append row error"

	if recieved != expected {
		t.Errorf("incorrect error message, received: %s, expected: %s", recieved, expected)
	}
}
