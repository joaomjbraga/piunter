package core

import "testing"

func TestIsVerboseEnabled(t *testing.T) {
	if isVerboseEnabled(false, false) {
		t.Fatal("expected verbose to be disabled by default")
	}

	if !isVerboseEnabled(true, false) {
		t.Fatal("expected verbose to be enabled when requested")
	}

	if !isVerboseEnabled(true, true) {
		t.Fatal("expected verbose to be enabled when config is verbose")
	}
}
