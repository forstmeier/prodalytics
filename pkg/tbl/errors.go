package tbl

import "fmt"

const packageName = "tbl"

// AppendRowError wraps errors returned by helper method
// calls in the tbl.Tabler.AppendRow method.
type AppendRowError struct {
	err error
}

func (e *AppendRowError) Error() string {
	return fmt.Sprintf("[%s] [append row]: %s", packageName, e.err.Error())
}
