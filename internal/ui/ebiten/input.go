// File internal/ui/ebiten/input.go
package ebiten

import (
	"dvonn_go/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
)

type dragState struct {
	active bool
	from   game.Coordinate
}

var d dragState

var (
	// clickStep = 0: nothing selected
	// clickStep = 1: waiting for destination
	clickStep int

	fromCoord game.Coordinate
)

var (
	selected   bool
	selectedAt game.Coordinate
)

// pixelToCoord converts screen pixels to a board coordinate.
func pixelToCoord(x, y int) (game.Coordinate, bool) {
	relX := float64(x) - offsetX
	relY := float64(y) - offsetY

	r := relY / (triangleR * 3.0 / 2)
	q := (relX / (triangleR * math.Sqrt(3))) - (r / 2)

	candidate, ok := indexFromAxial(q, r)
	if !ok || !onBoard(candidate) {
		return game.Coordinate{}, false
	}

	cx, cy := coordToScreen(candidate)
	dx := cx - float64(x)
	dy := cy - float64(y)

	if dx*dx+dy*dy <= (triangleR*0.5)*(triangleR*0.5) {
		return candidate, true
	}

	return game.Coordinate{}, false
}

// handleInput returns a user move (if any) based on mouse clicks.
func handleInput(gs *game.GameState) game.Move {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return nil
	}

	x, y := ebiten.CursorPosition()
	c, ok := pixelToCoord(x, y)
	if !ok {
		clickStep = 0
		selected = false
		return nil
	}

	if gs.Phase == game.Phase1 {
		game.RunPlacementPhase(gs, c.X, c.Y)
		return nil
	}

	if gs.Phase == game.Phase2 {
		switch clickStep {
		case 0:
			if !isMovable(gs, c) {
				return nil
			}
			fromCoord = c
			clickStep = 1
			selected = true
			selectedAt = c
			return nil
		case 1:
			for _, d := range destinations(gs, fromCoord) {
				if d == c {
					enterPerf()
					pl := game.TurnStateToPlayer(gs.Turn)
					clickStep = 0
					selected = false
					return game.JumpMove{Player: pl, From: fromCoord, To: c}
				}
			}
			clickStep = 0
			selected = false
			return nil
		}
	}

	return nil
}
