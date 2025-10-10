package common

import (
	"fmt"
	"strings"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

// ColorInfo returns a colored info symbol
func ColorInfo(text string) string {
	return fmt.Sprintf("%s%s%s", colorBlue, text, colorReset)
}

// ColorCyan returns text in cyan color
func ColorCyan(text string) string {
	return fmt.Sprintf("%s%s%s", colorCyan, text, colorReset)
}

// ColorSuccess returns text in green color
func ColorSuccess(text string) string {
	return fmt.Sprintf("%s%s%s", colorGreen, text, colorReset)
}

// ColorWarn returns text in yellow color
func ColorWarn(text string) string {
	return fmt.Sprintf("%s%s%s", colorYellow, text, colorReset)
}

// Separator returns a separator line of specified length
func Separator(length int) string {
	return strings.Repeat("-", length)
}
