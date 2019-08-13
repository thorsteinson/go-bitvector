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
			t.Errorf("Failed to add value %q", n)
		}

		v.Add(n)
		if !v.Contains(n) {
			t.Errorf("Failed idempotence test adding %q", n)
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
			t.Errorf("Failed to remove value %q", n)
		}

		v.Remove(n)
		if v.Contains(n) {
			t.Errorf("Failed idempotence test removing %q", n)
		}
	}
}

func Test_Values(t *testing.T) {
	bv := MakeVector(10)

	bv.Add(1)
	bv.Add(2)
	bv.Add(3)

	resultSlice := []int{1, 2, 3}

	equal := true

	for i, x := range bv.Values() {
		if x != resultSlice[i] {
			equal = false
			break
		}
	}

	if !equal {
		t.Errorf("Values doesn't have proper slice")
	}
}
