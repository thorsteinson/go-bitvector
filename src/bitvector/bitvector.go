package bitvector

const wordSize = 64

type Bitvec struct {
	words    []uint64
	capacity uint64
	size     int
}

func New(capacity int) *Bitvec {
	if capacity < 0 {
		panic("Negative capacity provided for Bitvec conctruction")
	}
	num_words := uint64(capacity)/wordSize + 1
	words := make([]uint64, num_words)

	bv := Bitvec{words: words, capacity: uint64(capacity)}
	return &bv
}

func (bv *Bitvec) Add(n int) {
	checkOOB(bv, n)
	if !(*bv).Contains(n) {
		wordIdx, innerIdx := index(n)
		(*bv).words[wordIdx] |= 1 << innerIdx
		(*bv).size++
	}
}

func (bv *Bitvec) Remove(n int) {
	checkOOB(bv, n)
	if (*bv).Contains(n) {
		wordIdx, innerIdx := index(n)
		(*bv).words[wordIdx] &= ^(1 << innerIdx)
		(*bv).size--
	}
}

func (bv *Bitvec) Values() (vals []int) {
	vals = []int{}

	var acc int

	for _, word := range bv.words {
		// We can skip the entire word if no bits are set
		if word == 0 {
			acc += wordSize
		} else {
			var bitIdx uint64
			for bitIdx = 0; bitIdx < wordSize; bitIdx++ {
				// Create a bitstring to test if the given bit is marked
				var bitstring uint64
				bitstring = 1 << bitIdx
				if word&bitstring > 0 {
					vals = append(vals, acc)
				}

				acc++
			}
		}
	}

	return vals
}

func (bv *Bitvec) Contains(n int) bool {
	checkOOB(bv, n)
	wordIdx, innerIdx := index(n)
	return (1<<innerIdx)&(*bv).words[wordIdx] > 0
}

func index(n int) (uint64, uint64) {
	m := uint64(n)
	wordIdx := m / wordSize
	innerIdx := m % wordSize
	return wordIdx, innerIdx
}

func checkOOB(bv *Bitvec, n int) {
	m := uint64(n)
	if m >= (*bv).capacity || m < 0 {
		panic("Out of index error")
	}
}

func (bv *Bitvec) Capacity() int { return int((*bv).capacity) }

func (bv *Bitvec) Size() int { return (*bv).size }

func (v1 *Bitvec) UnionWith(v2 *Bitvec) {
	capacityCheck(v1, v2)
	for i, word := range v2.words {
		v1.words[i] |= word
	}
}

func Union (v1 *Bitvec, v2 *Bitvec) *Bitvec {
	capacityCheck(v1, v2)
	v3 := New(v1.Capacity())
	for i, word := range v2.words {
		v3.words[i] = v1.words[i] | word
	}
	return v3
}

func capacityCheck(v1 *Bitvec, v2 *Bitvec) {
	if v1.Capacity() != v2.Capacity() {
		panic("Cannot union vectors of different capacity")
	}
}
