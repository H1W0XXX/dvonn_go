// File internal/ai/eval.go
package ai

import "dvonn_go/internal/game"

// Evaluate 静态评估：
// 1) 吃红子（Source piece）加分；
// 2) 吃对手的棋子更高分；
// 3) 离最近红子越近，加分越多；
// 4) 如果棋盘有一半以上空格，则根据控制棋子数量差评估。
func Evaluate(b *game.Board, me game.Player) int {
	// ——权重参数，可根据实测微调——
	const (
		wRedCapture    = 8  // 每吃到一个红子得分
		wEnemyCapture  = 3  // 每吃到一个对手棋子得分
		wProximityUnit = 4  // 离红子每近 1 格加分
		wControlPower  = 10 // 控制力差的权重
	)

	// 确定「我方」和「对手」的颜色
	var myCol, opCol game.Piece
	if me == game.PWhite {
		myCol, opCol = game.White, game.Black
	} else {
		myCol, opCol = game.Black, game.White
	}

	// 拿到所有红子（Source）的位置，以及最大可能距离，用来归一化
	sources := b.GetSourceCoordinates() // 需要在 Board 上实现这个方法
	maxDist := b.BoardDiameter()        // 得到棋盘直径（最大六边形距离）

	// 计算空位置的数量
	emptyCount := 0
	totalCount := len(b.Cells)
	for _, st := range b.Cells {
		if len(st) == 0 {
			emptyCount++
		}
	}

	// 计算控制棋子数量
	var myControl, opControl int
	for _, st := range b.Cells {
		if len(st) == 0 {
			continue
		}
		owner := st[0] // 栈顶决定控制方
		if owner == myCol {
			myControl++
		} else if owner == opCol {
			opControl++
		}
	}

	// 计算评分
	score := 0

	// 1) 吃红子 & 吃对手棋子加分
	for _, st := range b.Cells {
		if len(st) == 0 {
			continue
		}
		for _, p := range st {
			switch p {
			case game.Red:
				score += wRedCapture
			case opCol:
				score += wEnemyCapture
			}
		}
	}

	// 2) 离红子加分
	for coord, st := range b.Cells {
		if len(st) == 0 {
			continue
		}
		// 找到最小距离
		minD := maxDist
		for _, src := range sources {
			d := game.HexDistance(coord, src)
			if d < minD {
				minD = d
			}
		}
		score += (maxDist - minD) * wProximityUnit
	}

	// 3) 如果空位置占比超过一半，计算控制力差
	if float64(emptyCount)/float64(totalCount) > 0.5 {
		// 控制力差值
		controlDiff := myControl - opControl
		score += controlDiff * wControlPower
	}

	return score
}
