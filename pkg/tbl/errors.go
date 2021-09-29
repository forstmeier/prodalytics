package tbl

import "fmt"

const packageName = "tbl"

// ErrorAppendRow wraps errors returned by helper method
// calls in the tbl.Tabler.AppendRow method.
type ErrorAppendRow struct {
	err error
}

func (e *ErrorAppendRow) Error() string {
	return fmt.Sprintf("[%s] [append row]: %s", packageName, e.err.Error())
}
