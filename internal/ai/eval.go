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
	sources := b.GetSourceCoordinates()
	maxDist := b.BoardDiameter()
	if maxDist == 0 {
		maxDist = 1
	}

	totalCount := len(b.Cells)
	emptyCount := 0
	var myControl, opControl int
	score := 0

	for coord, st := range b.Cells {
		if len(st) == 0 {
			emptyCount++
			continue
		}

		owner := st[0] // 栈顶决定控制权
		ownerFactor := 0
		switch owner {
		case myCol:
			ownerFactor = 1
			myControl++
		case opCol:
			ownerFactor = -1
			opControl++
		}

		if ownerFactor == 0 {
			continue
		}

		// 1) 按控制方加权红子与被俘敌子的价值
		for _, p := range st {
			switch p {
			case game.Red:
				score += ownerFactor * wRedCapture
			case myCol:
				if owner == opCol {
					score -= wEnemyCapture // 我的子被对方控制，扣分
				}
			case opCol:
				if owner == myCol {
					score += wEnemyCapture // 对方的子被我控制，加分
				}
			}
		}

		// 2) 离最近红子的距离优势
		if len(sources) > 0 {
			minD := maxDist
			for _, src := range sources {
				d := game.HexDistance(coord, src)
				if d < minD {
					minD = d
				}
			}
			score += ownerFactor * (maxDist - minD) * wProximityUnit
		}
	}

	// 3) 空位占多数时，考虑控制力差
	if totalCount > 0 && float64(emptyCount)/float64(totalCount) > 0.5 {
		controlDiff := myControl - opControl
		score += controlDiff * wControlPower
	}

	return score
}
