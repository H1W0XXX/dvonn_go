// File internal/ai/search.go
package ai

import (
	"dvonn_go/internal/game"
	"math"
	"runtime"
	"sort"
	"sync"
)

// 并行搜索结果结构
type rootResult struct {
	mv    game.JumpMove
	score int
}

// SearchBestMove 入口：迭代加深 αβ，根节点多核并行，返回最佳 JumpMove
func SearchBestMove(gs *game.GameState, depth int) game.JumpMove {
	// 把 GOMAXPROCS 设为 CPU 核心数
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	me := game.TurnStateToPlayer(gs.Turn)

	// 生成所有有效的 JumpMove
	jmoves := getValidJumpMoves(gs, me)

	if len(jmoves) == 0 {
		// 如果没有有效移动，返回空移动
		return game.JumpMove{}
	}

	// 并行搜每个根走法
	results := make(chan rootResult, len(jmoves))
	var wg sync.WaitGroup
	for _, mv := range jmoves {
		wg.Add(1)
		go func(m game.JumpMove) {
			defer wg.Done()
			// 每个 goroutine 用自己独立的 transposition table，避免加锁开销
			tt := NewTT()

			// 克隆局面并执行一步
			clone := *gs
			clone.Board = gs.Board.Clone()

			// 直接应用移动，不再调用 RunMovementPhase 避免错误打印
			clone.Board = game.Apply(m, &clone.Board)
			clone.Turn = game.GetNextTurn(&clone.Board, me)

			// 检查是否还有合法移动
			if !game.HasAnyLegalMoves(&clone.Board, clone.Turn) {
				clone.Turn = game.End
			}

			// 深度减 1，alpha-beta 搜索
			v, _ := alphabeta(&clone, depth-1,
				math.MinInt32+1, math.MaxInt32-1,
				me, tt)
			// 因为这里返回的是下一手得分，父节点还要翻转视角
			results <- rootResult{mv: m, score: -v}
		}(mv)
	}

	// 等待所有并行完成
	wg.Wait()
	close(results)

	// 汇总选最优
	best := rootResult{score: math.MinInt32}
	for r := range results {
		if r.score > best.score {
			best = r
		}
	}
	return best.mv
}

// getValidJumpMoves 获取所有有效的跳子移动
func getValidJumpMoves(gs *game.GameState, player game.Player) []game.JumpMove {
	moves := game.GetPossibleMoves(&gs.Board)
	var validJumpMoves []game.JumpMove

	for _, m := range moves {
		if jm, ok := m.(game.JumpMove); ok && jm.Player == player {
			// 在这里预先验证移动的有效性
			if game.ValidMove(&gs.Board, jm) {
				validJumpMoves = append(validJumpMoves, jm)
			}
		}
	}

	// 简单排序：先尝试吃子多的
	sort.Slice(validJumpMoves, func(i, j int) bool {
		return len(*gs.Board.Cells[validJumpMoves[i].To.X][validJumpMoves[i].To.Y]) > len(*gs.Board.Cells[validJumpMoves[j].To.X][validJumpMoves[j].To.Y])
	})

	return validJumpMoves
}

// alphabeta 使用 TT
func alphabeta(gs *game.GameState, depth, alpha, beta int, me game.Player, tt *TT) (int, game.Move) {
	hash := Hash(&gs.Board, game.TurnStateToPlayer(gs.Turn))
	if v, mv, ok := tt.Lookup(hash, depth); ok {
		return v, mv
	}

	// 叶节点
	if depth == 0 || gs.Turn == game.End {
		v := Evaluate(&gs.Board, me)
		tt.Save(hash, depth, v, nil)
		return v, nil
	}

	// 获取当前玩家
	currentPlayer := game.TurnStateToPlayer(gs.Turn)

	// 生成所有有效的跳子移动
	jmoves := getValidJumpMoves(gs, currentPlayer)

	if len(jmoves) == 0 {
		// 没有有效移动，直接评估
		v := Evaluate(&gs.Board, me)
		tt.Save(hash, depth, v, nil)
		return v, nil
	}

	var bestMove game.Move
	for _, mv := range jmoves {
		// 复制 GameState
		clone := *gs
		clone.Board = gs.Board.Clone()

		// 直接应用移动
		clone.Board = game.Apply(mv, &clone.Board)
		clone.Turn = game.GetNextTurn(&clone.Board, currentPlayer)

		// 检查是否还有合法移动
		if !game.HasAnyLegalMoves(&clone.Board, clone.Turn) {
			clone.Turn = game.End
		}

		val, _ := alphabeta(&clone, depth-1, -beta, -alpha, me, tt)
		val = -val

		if val > alpha {
			alpha = val
			bestMove = mv
			if alpha >= beta {
				break // β剪枝
			}
		}
	}

	tt.Save(hash, depth, alpha, bestMove)
	return alpha, bestMove
}
