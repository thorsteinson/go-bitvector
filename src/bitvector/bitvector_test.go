package bitvector

import (
	"testing"
	"math/rand"
	"time"
)

const maxTestingSize = 10000
const testingRounds = 100

func TestAdd(t *testing.T) {
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	testSize := rand.Intn(maxTestingSize)
	bv := MakeVector(uint64(testSize))

	for i := 0; i < testingRounds; i++ {
		n := uint64(rand.Intn(testSize))
		bv.Add(n)
		if !bv.Contains(n) {
			t.Errorf("Value %q was not found after adding", n)
			t.Logf("Failed with seed: %q", seed)
		}
	}
}

// A really good test idea: We should add a ton of random values. We
// use the values function and ensure that all the values were in fact
// added. Then we call remove on all the same values and ensure that
// after every item has been removed, nothing is contained. This is a
// single comprehensive property test that can be run many times. Much
// better to use this, and then save the basic manual cases for our
// simpler tests

func TestRemove(t *testing.T) {
	bv := MakeVector(10)

	bv.Add(5)
	bv.Remove(5)

	if bv.Contains(5) {
		t.Errorf("Failed to remove 5")
	}
}

func Test_Values (t *testing.T) {
	bv := MakeVector(10)

	bv.Add(1)
	bv.Add(2)
	bv.Add(3)

	resultSlice := []uint64{1,2,3}

	equal := true

	for i,x := range bv.Values() {
		if x != resultSlice[i] {
			equal = false
			break
		}
	}

	if !equal {
		t.Errorf("Values doesn't have proper slice")
	}
}
