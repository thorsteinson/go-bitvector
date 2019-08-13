package bitvector

import (
	"testing"
)

const maxTestingSize = 10000
const testingRounds = 100

func Test_Add(t *testing.T) {
	v := MakeVector(1000)

	var tests = []int{0, 1, 2, 3, 4, 5, 6, 63, 64, 65, 127, 128, 129}

	for _, n := range tests {
		v.Add(n)
		if !v.Contains(n) {
			t.Errorf("Failed to add value %d", n)
		}

		v.Add(n)
		if !v.Contains(n) {
			t.Errorf("Failed idempotence test adding %d", n)
		}
	}
}

func TestRemove(t *testing.T) {
	v := MakeVector(1000)

	var tests = []int{0, 1, 2, 3, 4, 5, 6, 63, 64, 65, 127, 128, 129}

	for _, n := range tests {
		v.Add(n)
		v.Remove(n)
		if v.Contains(n) {
			t.Errorf("Failed to remove value %d", n)
		}

		v.Remove(n)
		if v.Contains(n) {
			t.Errorf("Failed idempotence test removing %d", n)
		}
	}
}

func Test_Values(t *testing.T) {
	v := MakeVector(1000)

	var tests = []int{0, 1, 2, 3, 4, 5, 6, 63, 64, 65, 127, 128, 129}

	for _, n := range tests {
		v.Add(n)
	}

	vals := v.Values()
	if len(vals) != len(tests) {
		t.Error("Value slice length differs from test length")
		t.Errorf("Expected value: %d", len(tests))
		t.Errorf("Found value: %d", len(vals))
	}

	for i, n := range vals {
		if n != tests[i] {
			t.Error("Unexpected value found while comparing values")
			t.Errorf("Expected value: %d", tests[i])
			t.Errorf("Found value: %d", n)
		}
	}
}

// TestOOB ensures that the bitvector will panic if we give it bad
// inputs that go outside the proper size.
func Test_OOB(t *testing.T) {
	v := MakeVector(100)

	negativeSizePanics := willPanic(func() { MakeVector(-1) })
	negativeValuePanics := willPanic(func() { v.Add(-1) })
	beyondSizePanics := willPanic(func() { v.Add(1000) })

	if !negativeSizePanics {
		t.Error("Failed to panic when creating with negative size")
	}
	if !negativeValuePanics {
		t.Error("Failed to panic when adding a negative value")
	}
	if !beyondSizePanics {
		t.Error("Failed to panic when adding beyond size of vector")
	}
}

// willPanic runs the provided function f and returns whether it
// panics or not.
func willPanic(f func()) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()
	f()
	return panicked
}
