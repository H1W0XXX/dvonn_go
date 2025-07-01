// File internal/game/move.go
package game

/*
本文件对应原 Haskell `Move` 模块，提供：
  - 走子合法性判定
  - 可能走子生成
  - 解析玩家输入到 Move
  - 轮换下一手
*/

// -----------------------------------------------------------------------------
// 基础判定工具
// -----------------------------------------------------------------------------

// isLinear 仅允许水平 / 垂直 / 45° 正斜线
func isLinear(m JumpMove) bool {
	dx := m.To.X - m.From.X
	dy := m.To.Y - m.From.Y
	return dx == 0 || dy == 0 || dx == dy
}

// isOnBoard 起终点均为棋盘合法坐标
func isOnBoard(b *Board, m JumpMove) bool {
	return validCoordinate(b, m.From) && validCoordinate(b, m.To)
}

// 坐标版
func coordOnBoard(b *Board, c Coordinate) bool {
	return validCoordinate(b, c)
}

// distance = max(|dx|, |dy|)
func hexDistance(m JumpMove) int {
	x1, y1 := m.From.X, m.From.Y
	x2, y2 := m.To.X, m.To.Y

	// 这是标准的 axial 坐标系六边形距离公式
	// 基于你提供的相邻关系，这应该是正确的
	dx := x2 - x1
	dy := y2 - y1

	// 六边形距离 = (|dx| + |dy| + |dx + dy|) / 2
	return (abs(dx) + abs(dy) + abs(dx+dy)) / 2
}

// playerOwns 看堆顶颜色
func playerOwns(p Player, st Stack) bool {
	if len(st) == 0 {
		return false
	}
	if p == PBlack {
		return st[0] == Black
	}
	return st[0] == White
}

// -----------------------------------------------------------------------------
// 合法性检查
// -----------------------------------------------------------------------------
func ValidMove(b *Board, mv Move) bool {
	switch m := mv.(type) {
	case JumpMove:
		// 1. 拿到起点堆栈和高度
		st := innerstack(b, m.From)
		h := len(st)

		// 2. 没棋子、高度不符、自身所有权、被包围，都直接 false
		if h == 0 ||
			!playerOwns(m.Player, st) ||
			isSurrounded(b, m.From) {
			return false
		}

		// 3. 定义轴坐标下的 6 个方向向量
		dirs := []Coordinate{
			{+1, -1}, // 右下
			{+1, 0},  // 右
			{0, +1},  // 右上
			{-1, +1}, // 左上
			{-1, 0},  // 左
			{0, -1},  // 左下
		}

		// 4. 看 m.To 是否正好等于 From + dir*h
		validDir := false
		for _, d := range dirs {
			if m.To.X == m.From.X+d.X*h &&
				m.To.Y == m.From.Y+d.Y*h {
				validDir = true
				break
			}
		}
		if !validDir {
			return false
		}

		// 5. 落点必须在棋盘内，且已有堆栈（不能落在空格上）
		if !coordOnBoard(b, m.To) ||
			!nonempty(b, m.To) {
			return false
		}

		return true

	case PlaceMove:
		// 放置阶段：用 At 而不是 Loc
		return validCoordinate(b, m.At) &&
			!nonempty(b, m.At)

	default:
		return false
	}
}

// -----------------------------------------------------------------------------
// 轮换逻辑
// -----------------------------------------------------------------------------

func getNextTurn(b *Board, justPlayed Player) TurnState {
	has := func(p Player) bool {
		for _, mv := range GetPossibleMoves(b) {
			if jm, ok := mv.(JumpMove); ok && jm.Player == p {
				return true
			}
		}
		return false
	}

	canB := has(PBlack)
	canW := has(PWhite)

	switch justPlayed {
	case PWhite:
		switch {
		case canB:
			return MoveBlack
		case canW:
			return MoveWhite
		default:
			return End
		}
	case PBlack:
		switch {
		case canW:
			return MoveWhite
		case canB:
			return MoveBlack
		default:
			return End
		}
	default:
		return End
	}
}

// -----------------------------------------------------------------------------
// 枚举所有可能跳子（Phase2）
// -----------------------------------------------------------------------------

func GetPossibleMoves(b *Board) []Move {
	var moves []Move
	cands := make([]Coordinate, 0, len(b.Cells))
	for c := range nonempties(b) {
		cands = append(cands, c)
	}
	for _, from := range cands {
		for _, to := range cands {
			for _, pl := range []Player{PWhite, PBlack} {
				mv := JumpMove{Player: pl, From: from, To: to}
				if ValidMove(b, mv) {
					moves = append(moves, mv)
				}
			}
		}
	}
	return moves
}

func HasAnyLegalMoves(b *Board, ts TurnState) bool {
	switch ts {
	case MoveWhite:
		for _, mv := range GetPossibleMoves(b) {
			if jm, ok := mv.(JumpMove); ok && jm.Player == PWhite {
				return true
			}
		}
	case MoveBlack:
		for _, mv := range GetPossibleMoves(b) {
			if jm, ok := mv.(JumpMove); ok && jm.Player == PBlack {
				return true
			}
		}
	}
	return false
}

// -----------------------------------------------------------------------------
// 内部小工具
// -----------------------------------------------------------------------------

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func max3(a, b, c int) int {
	if a >= b && a >= c {
		return a
	}
	if b >= a && b >= c {
		return b
	}
	return c
}
