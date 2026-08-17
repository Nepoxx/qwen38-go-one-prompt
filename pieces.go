package main

import "math/rand"

// Piece kinds. 0 means empty.
const (
	empty = iota
	I
	O
	T
	S
	Z
	J
	L
)

type cell struct{ x, y int }

// Base cells of each piece inside its bounding box.
var baseShapes = [8]struct {
	size  int
	cells []cell
}{
	1: {4, []cell{{0, 1}, {1, 1}, {2, 1}, {3, 1}}}, // I
	2: {2, []cell{{0, 0}, {1, 0}, {0, 1}, {1, 1}}}, // O
	3: {3, []cell{{1, 0}, {0, 1}, {1, 1}, {2, 1}}}, // T
	4: {3, []cell{{1, 0}, {2, 0}, {0, 1}, {1, 1}}}, // S
	5: {3, []cell{{0, 0}, {1, 0}, {1, 1}, {2, 1}}}, // Z
	6: {3, []cell{{0, 0}, {0, 1}, {1, 1}, {2, 1}}}, // J
	7: {3, []cell{{2, 0}, {0, 1}, {1, 1}, {2, 1}}}, // L
}

type piece struct {
	kind int
	rot  int // 0..3 clockwise rotations
	x, y int // board coords of the box top-left
}

// cells returns the piece's 4 cells after rot clockwise rotations.
func (p piece) cells() []cell {
	s := baseShapes[p.kind].size
	out := make([]cell, 4)
	for i, c := range baseShapes[p.kind].cells {
		x, y := c.x, c.y
		for r := 0; r < p.rot; r++ {
			x, y = s-1-y, x
		}
		out[i] = cell{x, y}
	}
	return out
}

func (p piece) rotated(dir int) piece {
	p.rot = (p.rot + dir + 4) % 4
	return p
}

// bag is a 7-bag randomizer: each piece appears once per bag.
type bag struct{ queue []int }

func (b *bag) next() int {
	if len(b.queue) == 0 {
		b.queue = []int{I, O, T, S, Z, J, L}
		rand.Shuffle(len(b.queue), func(i, j int) {
			b.queue[i], b.queue[j] = b.queue[j], b.queue[i]
		})
	}
	k := b.queue[0]
	b.queue = b.queue[1:]
	return k
}
