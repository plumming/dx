package table

import (
	"strings"
	"unicode/utf8"

	"github.com/plumming/dx/pkg/util"
)

const (
	alignCenter = 1
	alignLeft   = 2
)

func Pad(s, pad string, width int, align int) string {
	switch align {
	case alignCenter:
		return PadCenter(s, pad, width)
	case alignLeft:
		return PadLeft(s, pad, width)
	default:
		return PadRight(s, pad, width)
	}
}

func PadRight(s, pad string, width int) string {
	gap := width - utf8.RuneCountInString(util.Strip(s))
	if gap > 0 {
		return s + strings.Repeat(pad, gap)
	}
	return s
}

func PadLeft(s, pad string, width int) string {
	gap := width - utf8.RuneCountInString(util.Strip(s))
	if gap > 0 {
		return strings.Repeat(pad, gap) + s
	}
	return s
}

func PadCenter(s, pad string, width int) string {
	gap := width - utf8.RuneCountInString(util.Strip(s))
	if gap > 0 {
		gapLeft := int(float64(gap / 2))
		gapRight := gap - gapLeft
		return strings.Repeat(pad, gapLeft) + s + strings.Repeat(pad, gapRight)
	}
	return s
}
