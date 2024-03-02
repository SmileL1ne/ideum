package assert

import (
	"strings"
	"testing"
)

func Equal(t *testing.T, actual, expected interface{}) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; expected: %v", actual, expected)
	}
}

func StringContains(t *testing.T, actual, expectedStr string) {
	t.Helper()

	if !strings.Contains(actual, expectedStr) {
		t.Errorf("expected '%s' to contain '%s'", actual, expectedStr)
	}
}
