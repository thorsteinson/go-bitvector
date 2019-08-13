package bitvector

const wordSize = 64

type bitvec struct {
	words   []uint64
	maxSize uint64
	size    int
}

func MakeVector(size int) *bitvec {
	if size < 0 {
		panic("Negative size provided for Bitvec conctruction")
	}
	num_words := uint64(size)/wordSize + 1
	words := make([]uint64, num_words)

	bv := bitvec{words: words, maxSize: uint64(size)}
	return &bv
}

func (bv *bitvec) Add(n int) {
	checkOOB(bv, n)
	if !(*bv).Contains(n) {
		wordIdx, innerIdx := index(n)
		(*bv).words[wordIdx] |= 1 << innerIdx
		(*bv).size++
	}
}

func (bv *bitvec) Remove(n int) {
	checkOOB(bv, n)
	if (*bv).Contains(n) {
		wordIdx, innerIdx := index(n)
		(*bv).words[wordIdx] &= ^(1 << innerIdx)
		(*bv).size--
	}
}

func (bv *bitvec) Values() (vals []int) {
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

func (bv *bitvec) Contains(n int) bool {
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

func checkOOB(bv *bitvec, n int) {
	m := uint64(n)
	if m >= (*bv).maxSize || m < 0 {
		panic("Out of index error")
	}
}

func (bv *bitvec) Capacity() int { return int((*bv).maxSize) }

func (bv *bitvec) Size() int { return (*bv).size }
