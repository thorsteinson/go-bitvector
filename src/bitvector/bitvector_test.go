package bitvector

import (
	"math/rand"
	"testing"
	"time"
)

const maxTestingSize = 10000
const testingRounds = 100

func Test_Add(t *testing.T) {
	v := New(1000)

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
	v := New(1000)

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
	v := New(1000)

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
	v := New(100)

	negativeCapacityPanics := willPanic(func() { New(-1) })
	negativeValuePanics := willPanic(func() { v.Add(-1) })
	beyondCapacityPanics := willPanic(func() { v.Add(1000) })

	if !negativeCapacityPanics {
		t.Error("Failed to panic when creating with negative size")
	}
	if !negativeValuePanics {
		t.Error("Failed to panic when adding a negative value")
	}
	if !beyondCapacityPanics {
		t.Error("Failed to panic when adding beyond capacity of vector")
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

// TestComprehensive uses randomized trials to ensure that any
// combination of adding and removing elements leaves the bit vector
// in a proper state.
func TestComprehensive(t *testing.T) {
	const trials = 1000
	const maxCapacity = 500

	seed := time.Now().UnixNano()

	t.Logf("Setting seed as: %d", seed)
	rand.Seed(seed)

	for i := 0; i < trials; i++ {
		capacity := rand.Intn(maxCapacity)
		v := New(capacity)

		// Populate a slice random values that's half the capacity of the
		// selected capacity for the given trial
		sample := make([]int, capacity/2)
		for i := range sample {
			sample[i] = rand.Intn(capacity)
		}

		for _, n := range sample {
			v.Add(n)
		}

		// Ensure that every value from our sample was successfully
		// added
		for _, n := range sample {
			if !v.Contains(n) {
				t.Errorf("Missing added value: %d", n)
			}
		}

		for _, n := range sample {
			v.Remove(n)
		}

		// Ensure that every value was sucessfully removed
		for _, n := range sample {
			if v.Contains(n) {
				t.Errorf("Spurious value remaining: %d", n)
			}
		}
	}
}

func TestCapacity(t *testing.T) {
	caps := []int{1, 20, 64, 128, 300, 400, 500, 10000}

	for cap := range caps {
		v := New(cap)
		if v.Capacity() != cap {
			t.Errorf("Incorrect size. Expected %d => found %d", cap, v.Capacity())
		}
	}
}

func TestSize(t *testing.T) {
	v := New(100)
	sample := []int{1, 2, 3, 4, 5, 63, 64, 65}

	expectedSize := 0
	for n := range sample {
		v.Add(n)
		expectedSize++
		if v.Size() != expectedSize {
			t.Errorf("Size difference during adding. Expected %d => Found %d", expectedSize, v.Size())
		}
	}

	for n := range sample {
		v.Remove(n)
		expectedSize--
		if v.Size() != expectedSize {
			t.Errorf("Size difference during adding. Expected %d => Found %d", expectedSize, v.Size())
		}
	}
}

// CopyVec is a helper function that makes a copy of a vector, for
// testing and comparison later on.
func copyVec(v *Bitvec) *Bitvec {
	vCopy := New(v.Capacity())
	copy(vCopy.words, v.words)
	vCopy.capacity = v.capacity
	vCopy.size = v.size
	return vCopy
}

func Test_UnionWith(t *testing.T) {
	v1 := New(100)
	v2 := New(100)

	sample1 := []int{1, 2, 3, 4, 63, 64}
	sample2 := []int{65, 67, 88, 90}
	expectedVals := make([]int, len(sample1)+len(sample2))
	copy(expectedVals, sample1)
	copy(expectedVals[len(sample1):], sample2)

	for _, n := range sample1 {
		v1.Add(n)
	}

	for _, n := range sample2 {
		v2.Add(n)
	}
	v2Copy := copyVec(v2)

	t.Logf("v1.Values(): %v", v1.Values())
	t.Logf("v2.Values(): %v", v2.Values())

	v1.UnionWith(v2)

	// Checks that all of the elements from both bit vectors are
	// accounted for
	for i, val := range v1.Values() {
		if val != expectedVals[i] {
			t.Errorf("Unexpected value. Expected %d, Found %d", expectedVals[i], val)
		}
	}
	t.Logf("Values(): %v", v1.Values())
	t.Logf("Expected Values(): %v", expectedVals)

	// Checks that v2 has not been modified by the operation
	copiedVals := v2Copy.Values()
	for i, val := range v2.Values() {
		if val != copiedVals[i] {
			t.Errorf("Unexpected mutation in parameter vector. Expected value %d, Found %d", copiedVals[i], val)
		}
	}
}
