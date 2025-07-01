// File internal/ai/eval.go
package ai

import "dvonn_go/internal/game"

func Evaluate(b *game.Board, me game.Player) int {
	// 定义权重
	const (
		wOwn     = 1 // 自己原有棋子权重
		wCapture = 2 // 吃到对方棋子权重
	)
	// 计算我的得分和对方的得分
	var myScore, opScore int

	// 辅助：根据 player 拿到对应的颜色常量
	var myColor, opColor game.Piece
	if me == game.PWhite {
		myColor = game.White
		opColor = game.Black
	} else {
		myColor = game.Black
		opColor = game.White
	}

	for _, st := range b.Cells {
		if len(st) == 0 {
			continue
		}
		// st[0] 是栈顶颜色，决定谁“控制”这个堆
		owner := st[0]
		// 遍历这个堆里每一枚棋子，按颜色和权重累加分数
		switch owner {
		case myColor:
			for _, p := range st {
				if p == myColor {
					myScore += wOwn
				} else { // p == opColor
					myScore += wCapture
				}
			}
		case opColor:
			for _, p := range st {
				if p == opColor {
					opScore += wOwn
				} else { // p == myColor
					opScore += wCapture
				}
			}
		}
	}
	return myScore - opScore
}
