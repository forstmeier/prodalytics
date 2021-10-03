package evt

import "fmt"

const errorMessage = "package evt: %s"

// ConvertError wraps errors returned b helper
// method calls in the evt.Eventer.Convert method.
type ConvertError struct {
	err error
}

func (e *ConvertError) Error() string {
	return fmt.Sprintf(errorMessage, e.err.Error())
}
