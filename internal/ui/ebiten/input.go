// File internal/ui/ebiten/input.go
package ebiten

import (
	"dvonn_go/internal/game"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type dragState struct {
	active bool
	from   game.Coordinate
}

var d dragState

// 将鼠标像素 → 棋盘坐标；返回 (coord, 在棋盘内?)
// 将鼠标像素 → 棋盘坐标；返回 (coord, 在棋盘内?)
func pixelToCoord(x, y int) (game.Coordinate, bool) {
	// 将屏幕坐标转换为棋盘坐标
	// 考虑偏移量
	x -= offsetX
	y -= offsetY

	// 计算 r 值
	r := float64(y) / (triangleR * (3.0 / 2))

	// 计算 q 值
	q := (float64(x) / (triangleR * math.Sqrt(3))) - (math.Sqrt(3) / 2 * r)

	// 将 q 和 r 转为最近的整数坐标（可根据需要调整精度）
	roundedQ := int(math.Round(q))
	roundedR := int(math.Round(r))

	// 检查该坐标是否在棋盘内
	for c := range game.EmptyDvonn.Cells {
		cx, cy := coordToScreen(c)
		dx := cx - float64(x) - offsetX // 计算时需要重新考虑偏移
		dy := cy - float64(y) - offsetY
		if dx*dx+dy*dy < hexSize*hexSize/2 { // 圆心半径判定

			return game.Coordinate{X: roundedQ, Y: roundedR}, true
		}
	}

	return game.Coordinate{}, false
}

// 处理鼠标，返回生成的 Move（或 nil）
func handleInput(gs *game.GameState) game.Move {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		c, ok := pixelToCoord(x, y)
		fmt.Println(c)
		if !ok {
			return nil
		}
		if !d.active {
			// 鼠标按下，记录起点
			d = dragState{active: true, from: c}
		}
		return nil
	}
	// 鼠标释放
	if d.active {
		x, y := ebiten.CursorPosition()
		to, ok := pixelToCoord(x, y)

		d.active = false
		if !ok {
			return nil
		}
		// 放子阶段
		if gs.Phase == game.Phase1 {
			// 决定放哪种颜色
			var piece game.Piece
			switch gs.Turn {
			case game.PlacingRed:
				piece = game.Red
			case game.PlacingWhite:
				piece = game.White
			case game.PlacingBlack:
				piece = game.Black
			}
			return game.PlaceMove{Piece: piece, At: to}
		}
		// 跳子阶段
		pl := map[game.TurnState]game.Player{
			game.MoveWhite: game.PWhite,
			game.MoveBlack: game.PBlack,
		}[gs.Turn]
		return game.JumpMove{Player: pl, From: d.from, To: to}
	}
	return nil
}
