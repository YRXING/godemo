package gotest

import "testing"

func TestAdd(t *testing.T) {
	var a, b = 1, 2
	var expected = 3
	actual := Add(a, b)

	if actual != expected {
		t.Errorf("Add(%d %d) = %d;expected: %d", a, b, actual, expected)
	}
}
