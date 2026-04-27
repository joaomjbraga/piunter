package utils

import (
	"runtime"
	"strings"
)

func GetOptimalWorkers(taskCount int) int {
	cpuCount := runtime.NumCPU()
	workers := cpuCount
	if taskCount < workers {
		workers = taskCount
	}
	if workers < 1 {
		workers = 1
	}
	return workers
}

func SplitLines(s string) []string {
	if s == "" {
		return []string{}
	}
	var lines []string
	for _, line := range SplitSimple(s, '\n') {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func SplitColumns(s string) []string {
	return SplitSimple(s, '\t')
}

func SplitSimple(s string, sep rune) []string {
	var result []string
	var current string
	for _, r := range s {
		if r == sep {
			result = append(result, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	result = append(result, current)
	return result
}