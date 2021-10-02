package tbl

import "fmt"

const errorMessage = "package tbl: %s"

// AppendRowError wraps errors returned by helper method
// calls in the tbl.Tabler.AppendRow method.
type AppendRowError struct {
	err error
}

func (e *AppendRowError) Error() string {
	return fmt.Sprintf(errorMessage, e.err.Error())
}
