package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type state int

const (
	playing state = iota
	paused
	gameOver
)

// Gravity speed per level: 800ms at level 1, -70ms per level, floor 80ms.
func dropInterval(level int) int {
	ms := 800 - (level-1)*70
	if ms < 80 {
		ms = 80
	}
	return ms
}

// lineScores[n] is the base score for clearing n lines at once.
var lineScores = [5]int{0, 100, 300, 500, 800}

const (
	dasDelay = 12 // frames before held key repeats (~200ms)
	dasRate  = 3  // frames per repeat (~50ms)
	softMs   = 40 // gravity interval while soft-dropping
)

type game struct {
	board  board
	cur    piece
	next   int
	bag    bag
	score  int
	lines  int
	level  int
	state  state
	dropMs float64 // accumulated time since last gravity step
	dasDir int     // 0, -1 (left), +1 (right)
	dasT   int     // frames since the held key was pressed
}

func newGame() *game {
	g := &game{level: 1, state: playing}
	g.next = g.bag.next()
	g.spawn()
	return g
}

func (g *game) spawn() {
	k := g.next
	g.next = g.bag.next()
	s := baseShapes[k].size
	g.cur = piece{kind: k, x: (cols - s) / 2, y: 0}
	if g.board.collides(g.cur) {
		g.state = gameOver
	}
}

func (g *game) move(dx int) {
	p := g.cur
	p.x += dx
	if !g.board.collides(p) {
		g.cur = p
	}
}

// tryRotate applies a rotation with simple wall kicks.
func (g *game) tryRotate(dir int) {
	q := g.cur.rotated(dir)
	for _, k := range []cell{{0, 0}, {-1, 0}, {1, 0}, {0, -1}, {-2, 0}, {2, 0}, {-1, -1}, {1, -1}} {
		q2 := q
		q2.x += k.x
		q2.y += k.y
		if !g.board.collides(q2) {
			g.cur = q2
			return
		}
	}
}

// step advances gravity by one row, locking the piece if it cannot fall.
func (g *game) step(soft bool) {
	p := g.cur
	p.y++
	if !g.board.collides(p) {
		g.cur = p
		if soft {
			g.score++
		}
		return
	}
	g.lock()
}

func (g *game) hardDrop() {
	for {
		p := g.cur
		p.y++
		if g.board.collides(p) {
			break
		}
		g.cur = p
		g.score += 2
	}
	g.lock()
}

func (g *game) lock() {
	g.board.merge(g.cur)
	if n := g.board.clearLines(); n > 0 {
		g.score += lineScores[n] * g.level
		g.lines += n
		g.level = g.lines/10 + 1
	}
	g.dropMs = 0
	g.spawn()
}

// ghost returns where the current piece would land.
func (g *game) ghost() piece {
	p := g.cur
	for {
		q := p
		q.y++
		if g.board.collides(q) {
			return p
		}
		p = q
	}
}

// update runs one 60fps tick of input and gravity.
func (g *game) update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		*g = *newGame()
		return
	}
	switch g.state {
	case gameOver:
		return
	case paused:
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			g.state = playing
		}
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.state = paused
		return
	}

	// Horizontal movement with delayed auto shift.
	left := ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	right := ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD)
	justLeft := inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)
	justRight := inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)
	if justLeft {
		g.move(-1)
		g.dasDir, g.dasT = -1, 0
	} else if justRight {
		g.move(1)
		g.dasDir, g.dasT = 1, 0
	}
	if !left && !right {
		g.dasDir = 0
	}
	if g.dasDir != 0 {
		g.dasT++
		if g.dasT > dasDelay && (g.dasT-dasDelay)%dasRate == 0 {
			g.move(g.dasDir)
		}
	}

	soft := ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS)
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		g.tryRotate(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		g.tryRotate(-1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.hardDrop()
		return
	}

	interval := dropInterval(g.level)
	if soft && softMs < interval {
		interval = softMs
	}
	g.dropMs += 1000.0 / 60
	for g.dropMs >= float64(interval) {
		g.dropMs -= float64(interval)
		g.step(soft)
		if g.state != playing {
			return
		}
	}
}
