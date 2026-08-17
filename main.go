package main

import (
	"bytes"
	"image/color"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed fonts/DejaVuSans-Bold.ttf
var fontData []byte

// Layout constants.
const (
	px   = 28
	boardX = 20
	boardY = 20
	sideX  = boardX + cols*px + 20
	winW   = sideX + 170
	winH   = boardY + rows*px + 20
)

var (
	cBg     = color.RGBA{18, 20, 28, 255}
	cBoard  = color.RGBA{10, 12, 18, 255}
	cGrid   = color.RGBA{34, 38, 52, 255}
	cText   = color.RGBA{225, 228, 240, 255}
	cDim    = color.RGBA{120, 126, 148, 255}
	cAccent = color.RGBA{90, 200, 250, 255}
)

// Guideline colors for the seven pieces.
var pieceColor = [8]color.Color{
	I: color.RGBA{90, 220, 235, 255},
	O: color.RGBA{240, 200, 60, 255},
	T: color.RGBA{170, 90, 220, 255},
	S: color.RGBA{90, 200, 90, 255},
	Z: color.RGBA{230, 80, 80, 255},
	J: color.RGBA{80, 110, 230, 255},
	L: color.RGBA{240, 150, 60, 255},
}

var (
	faceTitle = newFace(26)
	faceBig   = newFace(20)
	faceMid   = newFace(14)
	faceSmall = newFace(12)
)

func newFace(size float64) *text.GoTextFace {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		panic(err)
	}
	return &text.GoTextFace{Source: s, Size: size}
}

type app struct {
	g  *game
	cb bool // colorblind mode: unique marker per piece kind
}

func newApp() *app { return &app{g: newGame()} }

func (a *app) Layout(w, h int) (int, int) { return winW, winH }

func (a *app) Update() error {
	a.g.update()
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		a.cb = !a.cb
	}
	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	screen.Fill(cBg)
	a.drawBoard(screen)
	a.drawSide(screen)
	if a.g.state == paused {
		a.overlay(screen, "PAUSED", "press P to resume")
	} else if a.g.state == gameOver {
		a.overlay(screen, "GAME OVER", "press R to restart")
	}
}

func (a *app) drawBoard(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, boardX, boardY, cols*px, rows*px, cBoard)
	for x := 0; x <= cols; x++ {
		ebitenutil.DrawRect(screen, boardX+float64(x)*px, boardY, 1, rows*px, cGrid)
	}
	for y := 0; y <= rows; y++ {
		ebitenutil.DrawRect(screen, boardX, boardY+float64(y)*px, cols*px, 1, cGrid)
	}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if k := a.g.board[y][x]; k != empty {
				a.blockAt(screen, boardX+float64(x)*px, boardY+float64(y)*px, pieceColor[k], k)
			}
		}
	}
	if a.g.state == playing {
		gh := a.g.ghost()
		for _, c := range gh.cells() {
			a.ghostBlock(screen, boardX+float64(gh.x+c.x)*px, boardY+float64(gh.y+c.y)*px, pieceColor[a.g.cur.kind])
		}
		for _, c := range a.g.cur.cells() {
			a.blockAt(screen, boardX+float64(a.g.cur.x+c.x)*px, boardY+float64(a.g.cur.y+c.y)*px, pieceColor[a.g.cur.kind], a.g.cur.kind)
		}
	}
}

// block draws a beveled cell, with a per-kind marker in colorblind mode.
func (a *app) block(screen *ebiten.Image, x, y float64, c color.Color) {
	a.blockAt(screen, x, y, c, 0)
}

func (a *app) blockAt(screen *ebiten.Image, x, y float64, c color.Color, kind int) {
	ebitenutil.DrawRect(screen, x, y, px, px, c)
	ebitenutil.DrawRect(screen, x, y, px, 4, lighten(c, 40))
	ebitenutil.DrawRect(screen, x, y, 4, px, lighten(c, 40))
	ebitenutil.DrawRect(screen, x, y+px-4, px, 4, darken(c, 45))
	ebitenutil.DrawRect(screen, x+px-4, y, 4, px, darken(c, 45))
	if a.cb && kind != 0 {
		a.marker(screen, x, y, c, kind)
	}
}

// marker draws a shape unique to each piece kind so pieces stay
// distinguishable without color.
func (a *app) marker(screen *ebiten.Image, x, y float64, c color.Color, kind int) {
	w := color.RGBA{255, 255, 255, 255}
	cx, cy := x+px/2, y+px/2
	switch kind {
	case I: // horizontal bar
		ebitenutil.DrawRect(screen, x+6, cy-2.5, px-12, 5, w)
	case S: // vertical bar
		ebitenutil.DrawRect(screen, cx-2.5, y+6, 5, px-12, w)
	case L: // square
		ebitenutil.DrawRect(screen, cx-5, cy-5, 10, 10, w)
	case J: // dot
		ebitenutil.DrawCircle(screen, cx, cy, 5.5, w)
	case O: // ring
		ebitenutil.DrawCircle(screen, cx, cy, 7, w)
		ebitenutil.DrawCircle(screen, cx, cy, 3.5, c)
	case Z: // X
		a.line(screen, x+7, y+7, x+px-7, y+px-7, w)
		a.line(screen, x+7, y+px-7, x+px-7, y+7, w)
	case T: // triangle
		a.line(screen, cx, cy-7, cx-8, cy+6, w)
		a.line(screen, cx-8, cy+6, cx+8, cy+6, w)
		a.line(screen, cx+8, cy+6, cx, cy-7, w)
	}
}

func (a *app) line(screen *ebiten.Image, x1, y1, x2, y2 float64, c color.Color) {
	ebitenutil.DrawLine(screen, x1, y1, x2, y2, c)
	ebitenutil.DrawLine(screen, x1+1, y1, x2+1, y2, c)
	ebitenutil.DrawLine(screen, x1, y1+1, x2, y2+1, c)
}

// ghostBlock draws a translucent landing preview.
func (a *app) ghostBlock(screen *ebiten.Image, x, y float64, c color.Color) {
	r, g, b, _ := c.RGBA()
	ebitenutil.DrawRect(screen, x+2, y+2, px-4, px-4, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 70})
}

func (a *app) drawSide(screen *ebiten.Image) {
	a.text(screen, "TETRIS", faceTitle, cAccent, sideX, boardY, false)
	a.text(screen, "SCORE", faceMid, cDim, sideX, boardY+52, false)
	a.text(screen, itoa(a.g.score), faceBig, cText, sideX, boardY+70, false)
	a.text(screen, "LEVEL", faceMid, cDim, sideX, boardY+112, false)
	a.text(screen, itoa(a.g.level), faceBig, cText, sideX, boardY+130, false)
	a.text(screen, "LINES", faceMid, cDim, sideX, boardY+172, false)
	a.text(screen, itoa(a.g.lines), faceBig, cText, sideX, boardY+190, false)

	a.text(screen, "NEXT", faceMid, cDim, sideX, boardY+232, false)
	a.drawNext(screen)

	a.text(screen, "← →  move", faceSmall, cDim, sideX, boardY+340, false)
	a.text(screen, "↑ / Z  rotate", faceSmall, cDim, sideX, boardY+360, false)
	a.text(screen, "X  rotate ccw", faceSmall, cDim, sideX, boardY+380, false)
	a.text(screen, "↓  soft drop", faceSmall, cDim, sideX, boardY+400, false)
	a.text(screen, "space  hard drop", faceSmall, cDim, sideX, boardY+420, false)
	a.text(screen, "P pause  R reset", faceSmall, cDim, sideX, boardY+440, false)
	a.text(screen, "C  colorblind mode", faceSmall, cDim, sideX, boardY+460, false)
	if a.cb {
		a.text(screen, "COLORBLIND: ON", faceMid, cAccent, sideX, boardY+484, false)
	} else {
		a.text(screen, "COLORBLIND: OFF", faceMid, cDim, sideX, boardY+484, false)
	}
}

// drawNext renders the upcoming piece centered in a 4x4 box.
func (a *app) drawNext(screen *ebiten.Image) {
	ox, oy := float64(sideX+20), float64(boardY+252)
	p := piece{kind: a.g.next}
	cells := p.cells()
	minX, minY, maxX, maxY := 9, 9, -1, -1
	for _, c := range cells {
		if c.x < minX {
			minX = c.x
		}
		if c.y < minY {
			minY = c.y
		}
		if c.x > maxX {
			maxX = c.x
		}
		if c.y > maxY {
			maxY = c.y
		}
	}
	w := float64(maxX-minX+1) * px
	h := float64(maxY-minY+1) * px
	cx, cy := ox+(4*px-w)/2, oy+(4*px-h)/2
	for _, c := range cells {
		a.blockAt(screen, cx+float64(c.x-minX)*px, cy+float64(c.y-minY)*px, pieceColor[a.g.next], a.g.next)
	}
}

func (a *app) overlay(screen *ebiten.Image, title, sub string) {
	ebitenutil.DrawRect(screen, boardX, boardY, cols*px, rows*px, color.RGBA{0, 0, 0, 160})
	a.text(screen, title, faceTitle, cText, boardX+cols*px/2, boardY+rows*px/2-30, true)
	a.text(screen, sub, faceMid, cDim, boardX+cols*px/2, boardY+rows*px/2+14, true)
}

// text draws s at (x,y); center horizontally when center is true.
func (a *app) text(screen *ebiten.Image, s string, f *text.GoTextFace, c color.Color, x, y float64, center bool) {
	if center {
		w, _ := text.Measure(s, f, f.Size)
		x -= w / 2
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale = scale(c)
	text.Draw(screen, s, f, op)
}

// lighten/darken shift an RGBA color by d (negative darkens).
func lighten(c color.Color, d int) color.Color {
	r, g, b, _ := c.RGBA()
	adj := func(v uint32) uint8 {
		x := int(v>>8) + d
		if x < 0 {
			x = 0
		}
		if x > 255 {
			x = 255
		}
		return uint8(x)
	}
	return color.RGBA{adj(r), adj(g), adj(b), 255}
}

func darken(c color.Color, d int) color.Color { return lighten(c, -d) }

func scale(c color.Color) ebiten.ColorScale {
	var cs ebiten.ColorScale
	cs.ScaleWithColor(c)
	return cs
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [12]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

func main() {
	ebiten.SetWindowSize(winW, winH)
	ebiten.SetWindowTitle("Tetris")
	if err := ebiten.RunGame(newApp()); err != nil {
		panic(err)
	}
}
