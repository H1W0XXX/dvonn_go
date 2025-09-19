package ebiten

import (
	"dvonn_go/internal/game"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math"
)

var (
	lineColor      = color.RGBA{0x44, 0x44, 0x44, 0xFF}
	circleColor    = color.RGBA{0xB0, 0xC4, 0xDE, 0xFF}
	highlightBlue  = color.RGBA{0x00, 0x66, 0xFF, 0xFF}
	highlightGreen = color.RGBA{0x00, 0xCC, 0x66, 0xFF}
)

var boardBG *ebiten.Image

const (
	triangleR = 70 //

	offsetX   = 650
	offsetY   = 360
	fillColor = 0xDDDDDDFF

	canvasW = 1300
	canvasH = 700
)

var triangleH = math.Sqrt(3) * triangleR / 2

func init() {

	boardBG = ebiten.NewImage(canvasW, canvasH)

	forEachCoordinate(func(c game.Coordinate) {
		drawCircle(boardBG, c)
		drawCoordinate(boardBG, c)
	})

	drawGridLines(boardBG)
}

func drawCoordinate(dst *ebiten.Image, c game.Coordinate) {
	cx, cy := coordToScreen(c)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("(%d,%d)", c.X, c.Y), int(cx-10), int(cy-10))
}

func drawCircle(dst *ebiten.Image, c game.Coordinate) {
	screenX, screenY := coordToScreen(c)
	radius := float64(triangleR) / 3
	ebitenutil.DrawCircle(dst, screenX, screenY, radius, circleColor)
}

func drawGridLines(dst *ebiten.Image) {
	forEachCoordinate(func(c1 game.Coordinate) {
		for _, c2 := range getAdjacentCoordinates(c1) {
			if !onBoard(c2) {
				continue
			}
			screenX1, screenY1 := coordToScreen(c1)
			screenX2, screenY2 := coordToScreen(c2)
			ebitenutil.DrawLine(dst, screenX1, screenY1, screenX2, screenY2, lineColor)
		}
	})
}

func getAdjacentCoordinates(c game.Coordinate) []game.Coordinate {
	var adjacent []game.Coordinate

	deltas := []game.Coordinate{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, -1}, {-1, 1},
	}

	for _, delta := range deltas {
		adjacent = append(adjacent, game.Coordinate{X: c.X + delta.X, Y: c.Y + delta.Y})
	}
	return adjacent
}

func drawCircleColored(dst *ebiten.Image, c game.Coordinate, col color.Color) {
	x, y := coordToScreen(c)
	r := float64(triangleR) / 3
	ebitenutil.DrawCircle(dst, x, y, r, col)
}
