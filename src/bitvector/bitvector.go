package bitvector

const wordSize = 64

type bitvec struct {
	words   []uint64
	maxSize uint64
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
	wordIdx, innerIdx := index(n)
	(*bv).words[wordIdx] |= 1 << innerIdx
}

func (bv *bitvec) Remove(n int) {
	checkOOB(bv, n)
	wordIdx, innerIdx := index(n)
	(*bv).words[wordIdx] &= ^(1 << innerIdx)
}

func (bv *bitvec) Values() (vals []int) {
	vals = []int{}

	var acc int

	for _, word := range bv.words {
		// Do a check here, if the word has a value of zero, that
		// means every bit is zero, and we can skip this entire
		// word. This can save us from a bunch of comparisons in a
		// more sparsely populated set.
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

	return vals
}

func (bv *bitvec) Contains(n int) bool {
	checkOOB(bv, n)
	wordIdx, innerIdx := index(n)
	return (1 << innerIdx) & (*bv).words[wordIdx] > 0
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
