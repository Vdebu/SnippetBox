package assert

import (
	"strings"
	"testing"
)

func StringContains(t *testing.T, actual, expectedSubString string) {
	t.Helper()
	// 判断后者是否在前者中
	if !strings.Contains(actual, expectedSubString) {
		t.Errorf("got: %q;expected to contain: %q", actual, expectedSubString)
	}
}

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v;want: %v", actual, expected)
	}
}
