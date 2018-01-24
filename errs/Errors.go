package errs

import (
	"fmt"
)

// Errorf shortcut to format errors
func Errorf(msg string, e error, args ...interface{}) string {
	msg = fmt.Sprintf(msg, args...)
	return msg + ": " + e.Error()
}
