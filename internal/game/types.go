// File internal/game/types.go
package game

// Piece 代表单枚棋子的颜色
type Piece uint8

const (
	Red Piece = iota
	White
	Black
)

// Player 用堆顶颜色判定所属方
type Player uint8

const (
	PWhite Player = iota
	PBlack
)

// Coordinate 采用整数轴向坐标；可直接作为 map key
type Coordinate struct {
	X, Y int
}

// Stack 为一列自顶向下的棋子
type Stack []Piece

// Board 持有棋盘上的所有坐标栈及弃子区
type Board struct {
	Cells   [BoardWidth][BoardHeight]*Stack // 空格用 nil / len==0 表示
	Discard Stack
}

// ---- Move 定义 -------------------------------------------------------------

type Move interface{ isMove() }

type JumpMove struct {
	Player   Player
	From, To Coordinate
}

func (JumpMove) isMove() {}

type PlaceMove struct {
	Piece Piece
	At    Coordinate
}

func (PlaceMove) isMove() {}

// PieceFromPlayer 把 Player 转 Piece（堆顶用）
func (p Player) Piece() Piece {
	if p == PWhite {
		return White
	}
	return Black
}
