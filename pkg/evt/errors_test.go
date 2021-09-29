package evt

import (
	"errors"
	"testing"
)

func TestErrorConvert(t *testing.T) {
	err := &ErrorConvert{
		err:      errors.New("mock convert error"),
		function: "function",
	}

	recieved := err.Error()
	expected := "[evt] [function]: mock convert error"

	if recieved != expected {
		t.Errorf("incorrect error message, received: %s, expected: %s", recieved, expected)
	}
}
