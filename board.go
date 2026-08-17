package main

const (
	cols = 10
	rows = 20
)

type board [rows][cols]int

// collides reports whether the piece overlaps the walls, floor, or locked cells.
// Cells above the top of the board are legal (spawn area).
func (b *board) collides(p piece) bool {
	for _, c := range p.cells() {
		x, y := p.x+c.x, p.y+c.y
		if x < 0 || x >= cols || y >= rows {
			return true
		}
		if y >= 0 && b[y][x] != empty {
			return true
		}
	}
	return false
}

// merge stamps the piece into the board.
func (b *board) merge(p piece) {
	for _, c := range p.cells() {
		y := p.y + c.y
		if y >= 0 {
			b[y][p.x+c.x] = p.kind
		}
	}
}

// clearLines removes full rows and shifts the rest down.
// Returns the number of lines cleared.
func (b *board) clearLines() int {
	cleared := 0
	write := rows - 1
	for y := rows - 1; y >= 0; y-- {
		full := true
		for x := 0; x < cols; x++ {
			if b[y][x] == empty {
				full = false
				break
			}
		}
		if full {
			cleared++
			continue
		}
		if write != y {
			b[write] = b[y]
		}
		write--
	}
	for ; write >= 0; write-- {
		b[write] = [cols]int{}
	}
	return cleared
}
