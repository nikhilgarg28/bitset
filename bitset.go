package bitset

const word = uint64(64)
const logword = uint(6)

type BitSet struct {
	length uint64
	bits   []uint64
}

func getSize(length uint64) uint64 {
	return uint64((length + word - 1) / word)
}

func New(length uint64) *BitSet {
	size := getSize(length)
	return &BitSet{length, make([]uint64, size)}
}

func getIndex(pos uint64) (q uint64, r uint) {
	q = pos >> logword
	r = uint(pos & (word - 1))
	return
}

func (b *BitSet) Length() uint64 {
	return b.length
}

func (b *BitSet) Get(pos uint64) bool {
	q, r := getIndex(pos)
	bit := (b.bits[q] >> r) & 1
	return bit != 0
}

func (b *BitSet) Set(pos uint64) bool {
	current := b.Get(pos)
	q, r := getIndex(pos)
	b.bits[q] |= (1 << r)
	return current
}

func (b *BitSet) Clear(pos uint64) bool {
	current := b.Get(pos)
	q, r := getIndex(pos)
	b.bits[q] &= ^(1 << r)
	return current
}

func (b *BitSet) Flip(pos uint64) bool {
	current := b.Get(pos)
	q, r := getIndex(pos)
	b.bits[q] ^= (1 << r)
	return current
}
