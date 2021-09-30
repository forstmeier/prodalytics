package evt

import "fmt"

const packageName = "evt"

// ConvertError wraps errors returned by functions and
// helper method calls in the evt.Eventer.Convert method.
type ConvertError struct {
	err      error
	function string
}

func (e *ConvertError) Error() string {
	return fmt.Sprintf("[%s] [%s]: %s", packageName, e.function, e.err.Error())
}
