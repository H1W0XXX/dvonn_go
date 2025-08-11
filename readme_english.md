# DVONN (Volcano) – Go Implementation

A desktop version of the classic DVONN board game, implemented in Go using the Ebiten engine.

## Introduction

This project recreates the DVONN (Volcano) game using Go and Ebiten. All image assets are sourced from [Board Game Arena – DVONN](https://www.boardgamearena.com/).

* **Game rules and state management**: `dvonn_go/internal/game`
* **GUI rendering and input handling**: `dvonn_go/internal/ui/ebiten`

Core logic is based on [gautammohan/dvonn](https://github.com/gautammohan/dvonn).

## Features

* Two game modes:

  * **PvE** (Player vs. AI)
  * **PvP** (Player vs. Player)
* Optional automatic placement of pieces during the initial setup (PvP only)
* 30 FPS frame rate limit

## Requirements

* Go 1.20 or higher
* Ebiten v2

## Build & Run

1. **Clone the repository** and navigate to the project root:

   ```bash
   git clone https://github.com/H1W0XXX/dvonn_go.git
   cd dvonn_go
   ```

2. **Build the executable**:

   ```bash
   go build -ldflags="-s -w" \
     -gcflags="all=-trimpath=${PWD}" \
     -asmflags="all=-trimpath=${PWD}" \
     -o dvonn.exe ./cmd/dvonn-gui/main.go
   ```

3. **Run the game**:

   ```bash
   ./dvonn.exe [flags]
   ```

## Command-Line Flags

| Flag    | Description                               | Default |
| ------- | ----------------------------------------- | ------- |
| `-auto` | Automatically place pieces in setup phase | `false` |
| `-mode` | Game mode: `pvp` or `pve`                 | `pve`   |

**Example:** Automatically place pieces in PvE mode

```bash
./dvonn.exe -mode=pve -auto
```

## Gameplay Overview

1. **Setup Phase**: Players take turns placing their pieces on empty spots until all pieces are placed.
2. **Movement Phase**: Players alternately move stacks. Each move distance equals the stack height, and moves must be in a straight line.
3. **End Condition**: When no more moves are possible, each player stacks all remaining stacks they control. The player with the tallest stack wins.

## References

* **Image Assets**: [Board Game Arena – DVONN](https://www.boardgamearena.com/)
* **Core Logic**: [gautammohan/dvonn](https://github.com/gautammohan/dvonn)

---

Contributions and feedback are welcome! Feel free to open issues or submit pull requests.
