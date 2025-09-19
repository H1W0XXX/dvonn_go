package ebiten

import (
	"dvonn_go/internal/game"
	"math"
)

func forEachCoordinate(fn func(game.Coordinate)) {
	game.ForEachPlayable(fn)
}

func onBoard(c game.Coordinate) bool { return game.IsPlayable(c) }

func axialFromIndex(c game.Coordinate) (float64, float64) {
	q, r := game.AxialFromIndex(c)
	return float64(q), float64(r)
}

func indexFromAxial(q, r float64) (game.Coordinate, bool) {
	iq := int(math.Round(q))
	ir := int(math.Round(r))
	coord, ok := game.IndexFromAxial(iq, ir)
	return coord, ok
}
