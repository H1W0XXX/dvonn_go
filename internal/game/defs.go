// File internal/game/defs.go
package game

import "fmt"

const (
	axialMinQ = -5
	axialMaxQ = 5
	axialMinR = -2
	axialMaxR = 2

	BoardWidth  = axialMaxQ - axialMinQ + 1
	BoardHeight = axialMaxR - axialMinR + 1
)

var (
	playableMask, playableCoords = buildPlayable()
)

// buildPlayable constructs the standard DVONN board mask (49 playable cells).
func buildPlayable() ([BoardWidth][BoardHeight]bool, []Coordinate) {
	var mask [BoardWidth][BoardHeight]bool
	coords := make([]Coordinate, 0, 49)

	for r := axialMinR; r <= axialMaxR; r++ {
		minQ, maxQ := rowBounds(r)
		for q := minQ; q <= maxQ; q++ {
			idx := axialToIndex(q, r)
			mask[idx.X][idx.Y] = true
			coords = append(coords, idx)
		}
	}

	return mask, coords
}

func rowBounds(r int) (int, int) {
	switch r {
	case -2:
		return -3, 5
	case -1:
		return -4, 5
	case 0:
		return -5, 5
	case 1:
		return -5, 4
	case 2:
		return -5, 3
	default:
		return 0, -1
	}
}

func axialToIndex(q, r int) Coordinate {
	return Coordinate{X: q - axialMinQ, Y: r - axialMinR}
}

func indexToAxial(c Coordinate) (int, int) {
	return c.X + axialMinQ, c.Y + axialMinR
}

// ForEachPlayable iterates over all playable coordinates.
func ForEachPlayable(fn func(Coordinate)) {
	for _, c := range playableCoords {
		fn(c)
	}
}

// IsPlayable reports whether the given coordinate is on the DVONN board.
func IsPlayable(c Coordinate) bool {
	return c.X >= 0 && c.X < BoardWidth && c.Y >= 0 && c.Y < BoardHeight && playableMask[c.X][c.Y]
}

// AxialFromIndex returns the axial (q, r) coordinate for a board index.
func AxialFromIndex(c Coordinate) (int, int) { return indexToAxial(c) }

// IndexFromAxial converts an axial coordinate to the internal board index.
func IndexFromAxial(q, r int) (Coordinate, bool) {
	if r < axialMinR || r > axialMaxR {
		return Coordinate{}, false
	}
	minQ, maxQ := rowBounds(r)
	if q < minQ || q > maxQ {
		return Coordinate{}, false
	}
	return axialToIndex(q, r), true
}

// -----------------------------------------------------------------------------
// 基础向量运算
// -----------------------------------------------------------------------------

// 人工实现简单向量加减（Go 不支持运算符重载）
func (c Coordinate) Add(o Coordinate) Coordinate { return Coordinate{c.X + o.X, c.Y + o.Y} }
func (c Coordinate) Sub(o Coordinate) Coordinate { return Coordinate{c.X - o.X, c.Y - o.Y} }
func (c Coordinate) Neg() Coordinate             { return Coordinate{-c.X, -c.Y} }
func (c Coordinate) String() string              { return fmt.Sprintf("(%d,%d)", c.X, c.Y) }

// Component = 连通块，由多个 Coordinate 组成的 set
type Component map[Coordinate]struct{}

// 空棋盘：按照棋盘形状初始化
func emptyBoard() Board {
	var b Board
	ForEachPlayable(func(c Coordinate) {
		b.Cells[c.X][c.Y] = nil
	})
	return b
}

// -----------------------------------------------------------------------------
// 默认棋盘状态（标准 DVONN 49 格棋盘）
// -----------------------------------------------------------------------------

var EmptyDvonn = emptyBoard()

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// EmptyMini ≈ 3×3 的简化棋盘（仅用于调试）
var EmptyMini = func() Board {
	var b Board
	for x := 1; x <= 3; x++ {
		for y := 1; y <= 3; y++ {
			c := Coordinate{X: x, Y: y}
			if IsPlayable(c) {
				b.Cells[c.X][c.Y] = nil
			}
		}
	}
	return b
}()

// -----------------------------------------------------------------------------
// 阶段 / 回合 / 错误
// -----------------------------------------------------------------------------

type TurnState string

const (
	PlacingRed   = "PlacingRed"
	PlacingWhite = "PlacingWhite"
	PlacingBlack = "PlacingBlack"
	MoveWhite    = "MoveWhite"
	MoveBlack    = "MoveBlack"
	Start        = "Start"
	End          = "End"
)

type GamePhase uint8

const (
	Phase1 GamePhase = iota // 摆子阶段
	Phase2                  // 行棋阶段
)

// -----------------------------------------------------------------------------
// Move / GameState / GameError
// -----------------------------------------------------------------------------

type GameState struct {
	Board     Board
	Turn      TurnState
	Phase     GamePhase
	PlaceStep int64
}

type GameError uint8

const (
	InvalidMove GameError = iota
	MoveParseError
)

func (e GameError) Error() string {
	switch e {
	case InvalidMove:
		return "invalid move"
	case MoveParseError:
		return "move parse error"
	default:
		return "unknown game error"
	}
}
