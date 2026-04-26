package utils

import (
	"strings"
)

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

type ParallelExecutor struct {
	maxWorkers int
}

func NewParallelExecutor(maxWorkers int) *ParallelExecutor {
	if maxWorkers <= 0 {
		maxWorkers = 4
	}
	return &ParallelExecutor{maxWorkers: maxWorkers}
}

type taskResult struct {
	index  int
	result any
	err    error
}

type voidTaskResult struct {
	index int
	err   error
}

func (e *ParallelExecutor) RunTasks(tasks []func() (any, error)) ([]any, []error) {
	results := make([]any, len(tasks))
	errors := make([]error, len(tasks))

	jobChan := make(chan int, len(tasks))
	resultChan := make(chan taskResult, e.maxWorkers)

	for i := range tasks {
		jobChan <- i
	}
	close(jobChan)

	for worker := 0; worker < e.maxWorkers; worker++ {
		go func() {
			for idx := range jobChan {
				result, err := tasks[idx]()
				resultChan <- taskResult{index: idx, result: result, err: err}
			}
		}()
	}

	for i := 0; i < len(tasks); i++ {
		res := <-resultChan
		results[res.index] = res.result
		errors[res.index] = res.err
	}

	return results, errors
}

func (e *ParallelExecutor) RunVoidTasks(tasks []func() error) []error {
	errors := make([]error, len(tasks))

	jobChan := make(chan int, len(tasks))
	resultChan := make(chan voidTaskResult, e.maxWorkers)

	for i := range tasks {
		jobChan <- i
	}
	close(jobChan)

	for worker := 0; worker < e.maxWorkers; worker++ {
		go func() {
			for idx := range jobChan {
				err := tasks[idx]()
				resultChan <- voidTaskResult{index: idx, err: err}
			}
		}()
	}

	for i := 0; i < len(tasks); i++ {
		res := <-resultChan
		errors[res.index] = res.err
	}

	return errors
}

type DirSizeResult struct {
	Path string
	Size int64
	Err  error
}

func GetDirSizesParallel(paths []string) []DirSizeResult {
	results := make([]DirSizeResult, len(paths))

	executor := NewParallelExecutor(4)
	hasErrors := false

	tasks := make([]func() (any, error), len(paths))
	for i, path := range paths {
		p := path
		tasks[i] = func() (any, error) {
			size, err := GetDirSize(p)
			return DirSizeResult{Path: p, Size: size, Err: err}, err
		}
	}

	_, errors := executor.RunTasks(tasks)

	for i, err := range errors {
		if err != nil {
			hasErrors = true
			results[i] = DirSizeResult{Path: paths[i], Size: 0, Err: err}
		}
	}

	if !hasErrors {
		for i, task := range tasks {
			result, _ := task()
			if dr, ok := result.(DirSizeResult); ok {
				results[i] = dr
			}
		}
	}

	return results
}

func AsyncGetDirSizes(paths []string, resultChan chan DirSizeResult, doneChan chan struct{}) {
	defer close(doneChan)

	workerCount := 4
	if len(paths) < workerCount {
		workerCount = len(paths)
	}

	var completed int

	for _, path := range paths {
		go func(p string) {
			size, err := GetDirSize(p)
			resultChan <- DirSizeResult{Path: p, Size: size, Err: err}
		}(path)
	}

	for completed < len(paths) {
		<-resultChan
		completed++
		if completed == len(paths) {
			break
		}
	}

	close(resultChan)
}