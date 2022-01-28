package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/xeonx/timeago"
)

func SafeTime(in *time.Time) string {
	if in == nil {
		return ""
	}
	t := *in
	return timeago.English.Format(t)
}

func TrimTemplate(instanceTemplate string, instanceType string) string {
	return strings.ReplaceAll(instanceTemplate, fmt.Sprintf("%s-instance-template-", instanceType), "")
}

func SafeIfAboveZero(number int) string {
	if number > 0 {
		return fmt.Sprintf("%d", number)
	}
	return ""
}
