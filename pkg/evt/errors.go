package evt

import "fmt"

const packageName = "evt"

// ErrorConvert wraps errors returned by functions and
// helper method calls in the evt.Eventer.Convert method.
type ErrorConvert struct {
	err      error
	function string
}

func (e *ErrorConvert) Error() string {
	return fmt.Sprintf("[%s] [%s]: %s", packageName, e.function, e.err.Error())
}
