package utils

import (
	"fmt"
)

var enableDebug = false

func SetDebug(enabled bool) {
	enableDebug = enabled
}

func Debug(msg string) {
	if enableDebug {
		fmt.Printf("\033[36m[DEBUG]\033[0m %s\n", msg)
	}
}

func Info(msg string) {
	fmt.Printf("  %s\n", msg)
}

func Error(msg string) {
	fmt.Printf("\033[31m  %s\033[0m\n", msg)
}

func Item(name, value string) {
	fmt.Printf("    \033[90m-\033[0m %-20s %s\n", name, value)
}

func Space() {
	fmt.Println()
}

func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

