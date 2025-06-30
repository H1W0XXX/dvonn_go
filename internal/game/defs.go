// File internal/game/defs.go
package game

import "fmt"

// -----------------------------------------------------------------------------
// 基本类型
// -----------------------------------------------------------------------------

// 人工实现常用向量运算（Go 不支持运算符重载）
func (c Coordinate) Add(o Coordinate) Coordinate { return Coordinate{c.X + o.X, c.Y + o.Y} }
func (c Coordinate) Sub(o Coordinate) Coordinate { return Coordinate{c.X - o.X, c.Y - o.Y} }
func (c Coordinate) Neg() Coordinate             { return Coordinate{-c.X, -c.Y} }
func (c Coordinate) String() string              { return fmt.Sprintf("(%d,%d)", c.X, c.Y) }

// Component = 连通块，由若干 Coordinate 组成的 set
type Component map[Coordinate]struct{}

// 空棋盘（传入合法坐标列表）
func emptyBoard(coords []Coordinate) Board {
	cells := make(map[Coordinate]Stack, len(coords))
	for _, c := range coords {
		cells[c] = nil // 空栈用 nil/len==0 表示
	}
	return Board{Cells: cells}
}

// -----------------------------------------------------------------------------
// 默认棋盘与测试棋盘
// -----------------------------------------------------------------------------

// EmptyDvonn ⟹ 49 格标准棋盘
var EmptyDvonn = func() Board {
	var coords []Coordinate

	// 手动定义六边形棋盘的所有坐标
	// 根据你提供的正确范围

	// 第1行：x从-3到5，y=-2 (9个位置)
	for x := -3; x <= 5; x++ {
		coords = append(coords, Coordinate{x, -2})
	}

	// 第2行：x从-4到5，y=-1 (10个位置)
	for x := -4; x <= 5; x++ {
		coords = append(coords, Coordinate{x, -1})
	}

	// 第3行：x从-5到5，y=0 (11个位置，最宽的一行)
	for x := -5; x <= 5; x++ {
		coords = append(coords, Coordinate{x, 0})
	}

	// 第4行：x从-5到4，y=1 (10个位置)
	for x := -5; x <= 4; x++ {
		coords = append(coords, Coordinate{x, 1})
	}

	// 第5行：x从-5到3，y=2 (9个位置)
	for x := -5; x <= 3; x++ {
		coords = append(coords, Coordinate{x, 2})
	}

	return emptyBoard(coords)
}()

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

// EmptyMini ⟹ 3×3 用于单元测试
var EmptyMini = func() Board {
	var coords []Coordinate
	for x := 1; x <= 3; x++ {
		for y := 1; y <= 3; y++ {
			coords = append(coords, Coordinate{x, y})
		}
	}
	return emptyBoard(coords)
}()

// -----------------------------------------------------------------------------
// 回合 / 阶段 / 玩家
// -----------------------------------------------------------------------------

// TurnState 表示当前轮到谁、在放子还是跳子
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
	Phase1 GamePhase = iota // 放子
	Phase2                  // 跳子
)

// -----------------------------------------------------------------------------
// Move / GameState / GameError
// -----------------------------------------------------------------------------

// 接口化便于 switch type
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
