//go:build cgo
// +build cgo

package state

import (
	"fmt"
	"time"
)

// formatDuration formats a duration (in seconds as int64 or time.Duration) to a human-readable string
func formatDuration(d interface{}) string {
	var duration time.Duration

	switch v := d.(type) {
	case int64:
		duration = time.Duration(v) * time.Second
	case time.Duration:
		duration = v
	default:
		return "-"
	}

	if duration == 0 {
		return "-"
	}
	if duration < time.Minute {
		return fmt.Sprintf("%ds", int(duration.Seconds()))
	}
	if duration < time.Hour {
		return fmt.Sprintf("%dm%ds", int(duration.Minutes()), int(duration.Seconds())%60)
	}
	return fmt.Sprintf("%dh%dm", int(duration.Hours()), int(duration.Minutes())%60)
}
