package main

import "testing"

func TestWallCollision(t *testing.T) {
	var b board
	if b.collides(piece{kind: O, x: 4, y: 5}) {
		t.Fatal("piece fully inside should not collide")
	}
	if !b.collides(piece{kind: O, x: -1, y: 5}) {
		t.Fatal("piece crossing left wall should collide")
	}
	if !b.collides(piece{kind: I, x: 7, y: 18}) {
		t.Fatal("I piece past right wall should collide")
	}
	if !b.collides(piece{kind: O, x: 4, y: rows - 1}) {
		t.Fatal("piece below floor should collide")
	}
}

func TestLineClear(t *testing.T) {
	var b board
	// Fill the bottom row except column 0.
	for x := 1; x < cols; x++ {
		b[rows-1][x] = T
	}
	// Rotated J fills the gap: cells land on (1,17),(0,17),(0,18),(0,19).
	p := piece{kind: J, rot: 1, x: -1, y: rows - 3}
	if b.collides(p) {
		t.Fatal("J should fit above the gap")
	}
	b.merge(p)
	if n := b.clearLines(); n != 1 {
		t.Fatalf("cleared %d lines, want 1", n)
	}
	if b[rows-1][0] != J || b[rows-2][0] != J || b[rows-2][1] != J {
		t.Fatal("J cells should rest on the floor after clear")
	}
	if b[rows-3][0] != empty || b[rows-3][1] != empty {
		t.Fatal("row above the J should be empty")
	}
}

func TestNoClearOnPartialRow(t *testing.T) {
	var b board
	for x := 0; x < cols-1; x++ {
		b[rows-1][x] = S
	}
	if n := b.clearLines(); n != 0 {
		t.Fatalf("cleared %d lines on a partial row, want 0", n)
	}
	if b[rows-1][0] != S {
		t.Fatal("partial row should be untouched")
	}
}
