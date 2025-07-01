// File internal/ui/ebiten/movehelper.go
package ebiten

import "dvonn_go/internal/game"

// getCurrentPlayer 将 TurnState 映射到对应的 Player
func getCurrentPlayer(gs *game.GameState) game.Player {
	switch gs.Turn {
	case game.MoveWhite, game.PlacingWhite:
		return game.PWhite
	case game.MoveBlack, game.PlacingBlack:
		return game.PBlack
	default:
		return game.PWhite
	}
}

// 返回当前玩家可移动的起点坐标列表
func movableCoords(gs *game.GameState) []game.Coordinate {
	seen := make(map[game.Coordinate]struct{})
	var res []game.Coordinate
	pl := getCurrentPlayer(gs)
	for _, mv := range game.GetPossibleMoves(&gs.Board) {
		jm, ok := mv.(game.JumpMove)
		if !ok || jm.Player != pl {
			continue
		}
		if _, exists := seen[jm.From]; !exists {
			seen[jm.From] = struct{}{}
			res = append(res, jm.From)
		}
	}
	return res
}

// 返回选中 from 后，所有合法落点坐标列表
func destinations(gs *game.GameState, from game.Coordinate) []game.Coordinate {
	var res []game.Coordinate
	pl := getCurrentPlayer(gs)
	for _, mv := range game.GetPossibleMoves(&gs.Board) {
		jm, ok := mv.(game.JumpMove)
		if !ok || jm.Player != pl || jm.From != from {
			continue
		}
		res = append(res, jm.To)
	}
	return res
}

// 判断指定坐标是否为当前玩家可移动的起点
func isMovable(gs *game.GameState, c game.Coordinate) bool {
	pl := getCurrentPlayer(gs)
	for _, mv := range game.GetPossibleMoves(&gs.Board) {
		jm, ok := mv.(game.JumpMove)
		if ok && jm.Player == pl && jm.From == c {
			return true
		}
	}
	return false
}
