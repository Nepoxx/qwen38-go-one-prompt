package main

import "testing"

func TestBagCoversAllPieces(t *testing.T) {
	var b bag
	seen := map[int]bool{}
	for i := 0; i < 7; i++ {
		seen[b.next()] = true
	}
	if len(seen) != 7 {
		t.Fatalf("first bag contained %d distinct pieces, want 7", len(seen))
	}
}

func TestBagRepeats(t *testing.T) {
	var b bag
	for i := 0; i < 14; i++ {
		b.next()
	}
	// After two full bags the queue must be refilled without panicking.
	if b.next() < I || b.next() > L {
		t.Fatal("piece kind out of range")
	}
}

func TestRotation(t *testing.T) {
	// T pointing down: (1,0),(0,1),(1,1),(2,1)
	// one clockwise rotation: (1,0),(1,1),(2,1),(1,2)
	want := map[cell]bool{
		{1, 0}: true, {1, 1}: true, {2, 1}: true, {1, 2}: true,
	}
	rotated := piece{kind: T, rot: 1}.cells()
	got := setOf(rotated)
	if len(got) != 4 || !sameCells(got, want) {
		t.Fatalf("T rot 1 = %v, want %v", got, want)
	}
	// Four rotations return to the start.
	if !sameCells(setOf(piece{kind: I, rot: 4}.cells()), setOf(piece{kind: I}.cells())) {
		t.Fatal("4 rotations did not return to base orientation")
	}
}

func sameCells(a, b map[cell]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for c := range a {
		if !b[c] {
			return false
		}
	}
	return true
}

func setOf(cs []cell) map[cell]bool {
	m := map[cell]bool{}
	for _, c := range cs {
		m[c] = true
	}
	return m
}
