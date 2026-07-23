package modules

import (
	"errors"
	"testing"
)

func TestIsPermissionError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{name: "permission denied", err: errors.New("permission denied"), want: true},
		{name: "access denied", err: errors.New("access denied"), want: true},
		{name: "operation not permitted", err: errors.New("operation not permitted"), want: true},
		{name: "other error", err: errors.New("file not found"), want: false},
		{name: "nil", err: nil, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isPermissionError(tc.err); got != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}
