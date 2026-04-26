package utils

import "github.com/joaomjbraga/piunter/pkg/types"

type CommandExecutor interface {
	Exec(command string, args ...string) types.CommandResult
	IsCommandAvailable(cmd string) bool
}

var defaultExecutor CommandExecutor = &realExecutor{}

type realExecutor struct{}

func (e *realExecutor) Exec(command string, args ...string) types.CommandResult {
	return Exec(command, args...)
}

func (e *realExecutor) IsCommandAvailable(cmd string) bool {
	return IsCommandAvailable(cmd)
}

func SetExecutor(exec CommandExecutor) {
	defaultExecutor = exec
}

func ResetExecutor() {
	defaultExecutor = &realExecutor{}
}

func GetExecutor() CommandExecutor {
	return defaultExecutor
}

type MockExecutor struct {
	execResults    map[string]types.CommandResult
	isAvailable  map[string]bool
	execCalls    []ExecCall
}

type ExecCall struct {
	Command string
	Args    []string
}

func NewMockExecutor() *MockExecutor {
	return &MockExecutor{
		execResults:   make(map[string]types.CommandResult),
		isAvailable: make(map[string]bool),
		execCalls:   make([]ExecCall, 0),
	}
}

func (m *MockExecutor) Exec(command string, args ...string) types.CommandResult {
	m.execCalls = append(m.execCalls, ExecCall{Command: command, Args: args})
	
	key := command
	for _, arg := range args {
		key += " " + arg
	}
	
	if result, ok := m.execResults[key]; ok {
		return result
	}
	
	if result, ok := m.execResults[command]; ok {
		return result
	}
	
	return types.CommandResult{
		Success: false,
		Stderr:  "mock: command not configured",
		Code:    1,
	}
}

func (m *MockExecutor) IsCommandAvailable(cmd string) bool {
	if available, ok := m.isAvailable[cmd]; ok {
		return available
	}
	return true
}

func (m *MockExecutor) WhenExec(command string, args ...string) *MockExecutor {
	key := command
	for _, arg := range args {
		key += " " + arg
	}
	m.execResults[key] = types.CommandResult{
		Success: true,
		Stdout:  "",
		Stderr:  "",
		Code:    0,
	}
	return m
}

func (m *MockExecutor) WhenExecResult(command string, result types.CommandResult) {
	m.execResults[command] = result
}

func (m *MockExecutor) WhenCommandAvailable(cmd string, available bool) {
	m.isAvailable[cmd] = available
}

func (m *MockExecutor) GetCalls() []ExecCall {
	return m.execCalls
}

func (m *MockExecutor) WasCalled(command string) bool {
	for _, call := range m.execCalls {
		if call.Command == command {
			return true
		}
	}
	return false
}

func (m *MockExecutor) CallCount(command string) int {
	count := 0
	for _, call := range m.execCalls {
		if call.Command == command {
			count++
		}
	}
	return count
}

func (m *MockExecutor) Reset() {
	m.execCalls = make([]ExecCall, 0)
}