package util

import (
	"fmt"
	"regexp"
)

var (
	InfoColor    = "\033[1;32m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	StatusColor  = "\033[1;34m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
)

// ColorInfo returns a new function that returns info-colorized (green) strings for the
// given arguments with fmt.Sprint().
func ColorInfo(in string) string {
	return fmt.Sprintf(InfoColor, in)
}

// ColorAnswer returns a new function that returns status-colorized (blue) strings for the
// given arguments with fmt.Sprint().
func ColorAnswer(in string) string {
	return fmt.Sprintf(StatusColor, in)
}

// ColorWarning returns a new function that returns warning-colorized (yellow) strings for the
// given arguments with fmt.Sprint().
func ColorWarning(in string) string {
	return fmt.Sprintf(WarningColor, in)
}

// ColorError returns a new function that returns error-colorized (red) strings for the
// given arguments with fmt.Sprint().
func ColorError(in string) string {
	return fmt.Sprintf(ErrorColor, in)
}

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

func Strip(str string) string {
	return re.ReplaceAllString(str, "")
}
