package assert

import (
	"strings"
	"testing"
)

// 用于检测字符串后者是否包含于前者
func StringContains(t *testing.T, actual, expectedSubString string) {
	t.Helper()
	// 判断后者是否在前者中
	if !strings.Contains(actual, expectedSubString) {
		t.Errorf("got: %q;expected to contain: %q", actual, expectedSubString)
	}
}

// 用于测试中各种类型的实际值比较
func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v;want: %v", actual, expected)
	}
}

// 判断错误是否是nil
func NilError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("got: %v;expected: nil", actual)
	}
}
