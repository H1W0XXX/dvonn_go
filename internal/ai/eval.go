// File internal/ai/eval.go
package ai

import "dvonn_go/internal/game"

func Evaluate(b *game.Board, me game.Player) int {
	// ——权重参数，可根据实测微调——
	const (
		wRedCapture    = 10 // 每吃到一个红子得分
		wEnemyCapture  = 2  // 每吃到一个对手棋子得分
		wProximityUnit = 4  // 离红子每近 1 格加分
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

	score := 0

	// 遍历所有格子
	for coord, st := range b.Cells {
		h := len(st)
		if h == 0 {
			continue
		}
		owner := st[0] // 栈顶决定控制方

		// 只关心我方控制的堆：其他堆对我没有直接正收益
		if owner != myCol {
			continue
		}

		// ——1) 吃红子 & 吃对手棋子加分——
		for _, p := range st {
			switch p {
			case game.Red:
				score += wRedCapture
			case opCol:
				score += wEnemyCapture
			}
		}

		// ——2) 距离最近红子 d 越小越好，加分 = (maxDist - d) * wProximityUnit——
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

	return score
}
