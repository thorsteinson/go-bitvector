package bitvector

const wordSize = 64

type bitvec struct {
	words   []uint64
	maxSize uint64
}

func MakeVector(size uint64) *bitvec {
	num_words := size/wordSize + 1
	words := make([]uint64, num_words)

	bv := bitvec{words: words, maxSize: size}
	return &bv
}

func (bv *bitvec) Add(n uint64) {
	checkOOB(bv, n)
	wordIdx, innerIdx := index(n)
	(*bv).words[wordIdx] |= 1 << innerIdx
}

func (bv *bitvec) Remove(n uint64) {
	checkOOB(bv, n)
	wordIdx, innerIdx := index(n)
	(*bv).words[wordIdx] &= ^(1 << innerIdx)
}

func (bv *bitvec) Values() (vals []uint64) {
	vals = []uint64{}

	var acc uint64

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

func (bv *bitvec) Contains(n uint64) bool {
	checkOOB(bv, n)
	wordIdx, innerIdx := index(n)
	return (1 << innerIdx) & (*bv).words[wordIdx] > 0
}

func index(n uint64) (uint64, uint64) {
	wordIdx := n / wordSize
	innerIdx := n % wordSize
	return wordIdx, innerIdx
}

func checkOOB(bv *bitvec, n uint64) {
	if n >= (*bv).maxSize {
		panic("Out of index error")
	}
}
