package util

import (
	"fmt"
	"regexp"

	"github.com/fatih/color"
)

var (
	InfoColor    = "\033[1;32m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	StatusColor  = "\033[1;34m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
)

var colorMap = map[string]color.Attribute{
	// formatting
	"bold":         color.Bold,
	"faint":        color.Faint,
	"italic":       color.Italic,
	"underline":    color.Underline,
	"blinkslow":    color.BlinkSlow,
	"blinkrapid":   color.BlinkRapid,
	"reversevideo": color.ReverseVideo,
	"concealed":    color.Concealed,
	"crossedout":   color.CrossedOut,

	// Foreground text colors
	"black":   color.FgBlack,
	"red":     color.FgRed,
	"green":   color.FgGreen,
	"yellow":  color.FgYellow,
	"blue":    color.FgBlue,
	"magenta": color.FgMagenta,
	"cyan":    color.FgCyan,
	"white":   color.FgWhite,

	// Foreground Hi-Intensity text colors
	"hiblack":   color.FgHiBlack,
	"hired":     color.FgHiRed,
	"higreen":   color.FgHiGreen,
	"hiyellow":  color.FgHiYellow,
	"hiblue":    color.FgHiBlue,
	"himagenta": color.FgHiMagenta,
	"hicyan":    color.FgHiCyan,
	"hiwhite":   color.FgHiWhite,

	// Background text colors
	"bgblack":   color.BgBlack,
	"bgred":     color.BgRed,
	"bggreen":   color.BgGreen,
	"bgyellow":  color.BgYellow,
	"BgBlue":    color.BgBlue,
	"bgmagenta": color.BgMagenta,
	"bgcyan":    color.BgCyan,
	"bgwhite":   color.BgWhite,

	// Background Hi-Intensity text colors
	"bghiblack":   color.BgHiBlack,
	"bghired":     color.BgHiRed,
	"bghigreen":   color.BgHiGreen,
	"bghiyellow":  color.BgHiYellow,
	"bghiblue":    color.BgHiBlue,
	"bghimagenta": color.BgHiMagenta,
	"bghicyan":    color.BgHiCyan,
	"bghiwhite":   color.BgHiWhite,
}

// GetColor returns the color for the list of colour names and option name.
func GetColor(optionName string, colorNames []string) (*color.Color, error) {
	var attributes []color.Attribute
	for _, colorName := range colorNames {
		a := colorMap[colorName]
		if a == color.Attribute(0) {
			return nil, fmt.Errorf("invalid color: " + optionName)
		}
		attributes = append(attributes, a)
	}
	return color.New(attributes...), nil
}

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
