package utils

import (
	"testing"

	"github.com/joaomjbraga/piunter/pkg/types"
)

func TestMockExecutor_BasicExec(t *testing.T) {
	mock := NewMockExecutor()

	mock.WhenExecResult("docker ps", types.CommandResult{
		Success: true,
		Stdout:  "container1\ncontainer2",
		Code:    0,
	})

	result := mock.Exec("docker", "ps")

	if !result.Success {
		t.Error("expected success")
	}

	if result.Stdout != "container1\ncontainer2" {
		t.Errorf("expected stdout, got %s", result.Stdout)
	}
}

func TestMockExecutor_NotConfigured(t *testing.T) {
	mock := NewMockExecutor()

	result := mock.Exec("docker", "ps")

	if result.Success {
		t.Error("expected failure for not configured command")
	}
}

func TestMockExecutor_WhenCommandAvailable(t *testing.T) {
	mock := NewMockExecutor()

	mock.WhenCommandAvailable("docker", true)
	mock.WhenCommandAvailable("npm", false)

	if !mock.IsCommandAvailable("docker") {
		t.Error("expected docker to be available")
	}

	if mock.IsCommandAvailable("npm") {
		t.Error("expected npm to not be available")
	}
}

func TestMockExecutor_GetCalls(t *testing.T) {
	mock := NewMockExecutor()

	mock.Exec("docker", "ps")
	mock.Exec("npm", "cache", "clean")
	mock.Exec("docker", "ps")

	calls := mock.GetCalls()

	if len(calls) != 3 {
		t.Errorf("expected 3 calls, got %d", len(calls))
	}

	if calls[0].Command != "docker" {
		t.Errorf("expected first command to be docker, got %s", calls[0].Command)
	}
}

func TestMockExecutor_WasCalled(t *testing.T) {
	mock := NewMockExecutor()

	mock.Exec("docker", "ps")
	mock.Exec("npm", "list")

	if !mock.WasCalled("docker") {
		t.Error("expected docker to be called")
	}

	if mock.WasCalled("podman") {
		t.Error("expected podman to not be called")
	}
}

func TestMockExecutor_CallCount(t *testing.T) {
	mock := NewMockExecutor()

	mock.Exec("docker", "ps")
	mock.Exec("docker", "ps")
	mock.Exec("npm", "list")

	if count := mock.CallCount("docker"); count != 2 {
		t.Errorf("expected 2 calls to docker, got %d", count)
	}
}

func TestMockExecutor_Reset(t *testing.T) {
	mock := NewMockExecutor()

	mock.Exec("docker", "ps")
	mock.Reset()

	if len(mock.GetCalls()) != 0 {
		t.Error("expected calls to be reset")
	}
}

func TestSetAndResetExecutor(t *testing.T) {
	mock := NewMockExecutor()

	SetExecutor(mock)

	if GetExecutor() != mock {
		t.Error("expected executor to be set")
	}

	ResetExecutor()

	if GetExecutor() == mock {
		t.Error("expected executor to be reset to real executor")
	}
}

func TestRealExecutor(t *testing.T) {
	executor := GetExecutor()

	result := executor.Exec("echo", "test")

	if !result.Success {
		t.Error("expected echo to succeed")
	}

	if result.Stdout != "test\n" {
		t.Errorf("expected 'test\\n', got %s", result.Stdout)
	}
}

func TestRealExecutor_IsCommandAvailable(t *testing.T) {
	executor := GetExecutor()

	if !executor.IsCommandAvailable("ls") {
		t.Error("expected ls to be available")
	}

	if executor.IsCommandAvailable("nonexistent_command_12345") {
		t.Error("expected nonexistent command to not be available")
	}
}