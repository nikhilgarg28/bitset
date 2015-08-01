package bitset

const word = uint64(64)
const logword = uint(6)

type Bitset struct {
	length uint64
	bits   []uint64
}

func getSize(length uint64) uint64 {
	return uint64((length + word - 1) / word)
}

func New(length uint64) *Bitset {
	size := getSize(length)
	return &Bitset{
		length,
		make([]uint64, size),
	}
}

func getIndex(pos uint64) (q uint64, r uint) {
	q = pos >> logword
	r = uint(pos & (word - 1))
	return
}

func (b *Bitset) Length() uint64 {
	return b.length
}

func (b *Bitset) Get(pos uint64) bool {
	q, r := getIndex(pos)
	bit := (b.bits[q] >> r) & 1
	return bit != 0
}

func (b *Bitset) Set(pos uint64) bool {
	current := b.Get(pos)
	q, r := getIndex(pos)
	b.bits[q] |= (1 << r)
	return current
}

func (b *Bitset) Clear(pos uint64) bool {
	current := b.Get(pos)
	q, r := getIndex(pos)
	b.bits[q] &= ^(1 << r)
	return current
}

func (b *Bitset) Flip(pos uint64) bool {
	current := b.Get(pos)
	q, r := getIndex(pos)
	b.bits[q] ^= (1 << r)
	return current
}
