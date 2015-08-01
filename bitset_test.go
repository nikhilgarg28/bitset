package bitset

import (
	"testing"
)

func TestInit(t *testing.T) {
	bs := New(1000)
	// initially all bits should be zero
	for i := uint64(0); i < 1000; i++ {
		if bs.Get(i) {
			t.Errorf("Bit %d is set when it should not be", i)
		}
	}
}

func TestGetSet(t *testing.T) {
	bs := New(1000)
	if bs.Set(0) == true {
		t.Error("Set on bit %d returned true when it should return false", 0)
	}

	if bs.Get(0) == false {
		t.Errorf("Bit %d is not set when it should be", 0)
	}

	if bs.Set(0) == false {
		t.Errorf("Set on bit %d returned false when it should return true", 0)
	}
}

func TestLargeSetGet(t *testing.T) {
	size := uint64(1) << 35
	bs := New(size)

	positions := []uint64{0, 1, 10, 1000, 1 << 32, size - 1}
	for _, position := range positions {
		if bs.Get(position) {
			t.Errorf("Bit %d should be false but is true", position)
		}

		bs.Set(position)

		if !bs.Get(position) {
			t.Errorf("Bit %d should be true but is false", position)
		}
	}
}

func TestLength(t *testing.T) {

	sizes := []uint64{0, 1, 10, 1000, 1 << 32, 1 << 33}
	for _, size := range sizes {
		bs := New(size)
		if bs.Length() != size {
			t.Errorf("Length should be %d", size)
		}
	}
}

func TestClear(t *testing.T) {
	bs := New(1000)
	pos := uint64(1)
	bs.Set(pos)

	if !bs.Get(pos) {
		t.Errorf("Bit %d should be set when it is not", pos)
	}

	if !bs.Clear(pos) {
		t.Error("Clearing bit %d should have returned true", pos)
	}

	if bs.Get(pos) {
		t.Errorf("Bit %d should have been cleared when it is not", pos)
	}
}

func TestFlip(t *testing.T) {
	bs := New(1000)
	pos := uint64(1)

	if bs.Get(pos) {
		t.Errorf("Bit %d should be clear when it is not", pos)
	}

	if bs.Flip(pos) {
		t.Error("Flipping bit %d should have returned false, but it did not", pos)
	}

	if !bs.Get(pos) {
		t.Errorf("Bit %d should have been set when it is not", pos)
	}
	if !bs.Flip(pos) {
		t.Error("Flipping bit %d should have returned true, but it did not", pos)
	}

	if bs.Get(pos) {
		t.Errorf("Bit %d should have been clear when it is not", pos)
	}
}
