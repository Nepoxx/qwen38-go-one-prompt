# Tetris

A complete Tetris game for the desktop, written in Go with [Ebiten](https://ebiten.org/) (OpenGL/OpenGL ES rendering).

## Build & run

Requires Go 1.26+ and the usual Ebiten system dependencies on Linux
(GCC, X11 development headers, ALSA).

```sh
go build -o tetris .
./tetris
```

Run the tests (game logic only, no display needed):

```sh
go test ./...
```

## Controls

| Key            | Action          |
|----------------|-----------------|
| ← → / A D     | Move            |
| ↑ / Z         | Rotate          |
| X             | Rotate ccw      |
| ↓ / S         | Soft drop (+1/cell) |
| Space         | Hard drop (+2/cell) |
| P             | Pause           |
| R             | Restart         |

## Gameplay

- Standard 10×20 playfield, all seven tetrominoes
- 7-bag randomizer (every piece appears once per bag)
- SRS-style rotation with wall kicks
- Ghost piece showing the landing position
- Scoring: 100/300/500/800 × level for 1–4 lines, plus drop bonuses
- Level up every 10 lines; gravity speeds up from 800 ms to 80 ms per row

## Files

- `pieces.go` — tetromino shapes, rotation, 7-bag randomizer
- `board.go` — grid, collision detection, line clearing
- `game.go` — game state, gravity, input (with DAS), scoring
- `main.go` — Ebiten window, rendering, HUD
- `fonts/` — DejaVu Sans Bold (free Bitstream Vera/DejaVu font license), embedded for the HUD text
