package utils

import (
	"fmt"
	"strings"
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

func Warn(msg string) {
	fmt.Printf("\033[33m  %s\033[0m\n", msg)
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

func List(items []struct {
	Name   string
	Value  string
	Success bool
}) {
	for _, item := range items {
		if item.Success {
			Item(item.Name, item.Value)
		} else {
			Item(item.Name, "\033[31merro\033[0m")
		}
	}
}

func ParseThreshold(value string) int {
	if value == "" {
		return 100
	}
	var threshold int
	fmt.Sscanf(value, "%d", &threshold)
	if threshold < 1 {
		return 1
	}
	if threshold > 10000 {
		return 10000
	}
	return threshold
}

func HasPrefixCI(s, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix))
}