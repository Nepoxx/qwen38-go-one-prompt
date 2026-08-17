# Tetris

Built from two prompts:

> Build a complete, polished Tetris game in Go.
>
> Make it a graphical desktop application, not a terminal application. Choose an appropriate Go game/windowing library and set up the project yourself.
>
> The game should feel like a real small game rather than a technical demo. It should have:
>
> - A standard 10×20 playfield
> - All seven tetrominoes
> - Piece movement and rotation
> - Correct collision detection
> - Gravity and progressively increasing difficulty
>
> This is a clean git repository. Also commit and push your changes. They should include a sensible gitignore and the commits should demonstrate progress

> Add support for colorblind users. It should be a game option.

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
| C             | Colorblind mode |

## Gameplay

- Standard 10×20 playfield, all seven tetrominoes
- 7-bag randomizer (every piece appears once per bag)
- SRS-style rotation with wall kicks
- Ghost piece showing the landing position
- Colorblind mode (C): every piece kind gets a unique white shape marker
  (bar, dot, ring, square, X, triangle) so pieces stay distinguishable
  without relying on color
- Scoring: 100/300/500/800 × level for 1–4 lines, plus drop bonuses
- Level up every 10 lines; gravity speeds up from 800 ms to 80 ms per row

## Files

- `pieces.go` — tetromino shapes, rotation, 7-bag randomizer
- `board.go` — grid, collision detection, line clearing
- `game.go` — game state, gravity, input (with DAS), scoring
- `main.go` — Ebiten window, rendering, HUD
- `fonts/` — DejaVu Sans Bold (free Bitstream Vera/DejaVu font license), embedded for the HUD text
