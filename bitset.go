package bitset

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

const word = uint64(64)
const logword = uint(6)
const bufsize = 1 << 16

// each channel write is 9 bytes, so if we were to write full disk sector
// of 64K bytes at once, we'd flush ~7112 writes
const channelsize = 7200

const (
	new = iota
	set
	clear
	flip
)

type command struct {
	code  uint8
	index uint64
}

type BitSet struct {
	length     uint64
	bits       []uint64
	writes     chan command
	binlogfile string
}

func getSize(length uint64) uint64 {
	return uint64((length + word - 1) / word)
}

func flush(filename string, writes chan command) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	// make a buffered writer
	w := bufio.NewWriterSize(f, bufsize)

	ticker := time.NewTimer(time.Second)

	buf := make([]byte, bufsize)
	long := make([]byte, 8)

	put := func(w *command) {
		binary.BigEndian.PutUint64(long, w.index)
		buf = append(buf, w.code)
		buf = append(buf, long...)
	}

	for t := range ticker.C {
		_ = t
		b := 0

	drain:
		for b+9 < bufsize {
			select {
			case write := <-writes:
				put(&write)
				b += 9
			default:
				break drain
			}
		}

		if b > 0 {
			// write a chunk
			if _, err := w.Write(buf[:b]); err != nil {
				panic(err)
			}

			if err = w.Flush(); err != nil {
				panic(err)
			}
		}
	}
}

func New(length uint64, binlogfile string) *BitSet {
	size := getSize(length)
	ret := &BitSet{
		length,
		make([]uint64, size),
		make(chan command, channelsize),
		binlogfile,
	}
	ret.writes <- command{new, length}
	go flush(binlogfile, ret.writes)
	return ret
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
	b.writes <- command{set, pos}
	return current
}

func (b *BitSet) Clear(pos uint64) bool {
	current := b.Get(pos)
	q, r := getIndex(pos)
	b.bits[q] &= ^(1 << r)
	b.writes <- command{clear, pos}
	return current
}

func (b *BitSet) Flip(pos uint64) bool {
	current := b.Get(pos)
	q, r := getIndex(pos)
	b.bits[q] ^= (1 << r)
	b.writes <- command{flip, pos}
	return current
}
