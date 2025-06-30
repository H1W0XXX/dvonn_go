// File internal/ui/ebiten/input.go
package ebiten

import (
	"dvonn_go/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type dragState struct {
	active bool
	from   game.Coordinate
}

var d dragState

// 将鼠标像素 → 棋盘坐标；返回 (coord, 在棋盘内
// 反算坐标
func pixelToCoord(x, y int) (game.Coordinate, bool) {
	// 将屏幕坐标转换为相对坐标
	relX := float64(x) - offsetX
	relY := float64(y) - offsetY

	// 计算 r 值
	r := relY / (triangleR * 3.0 / 2)

	// 计算 q 值
	q := (relX / (triangleR * math.Sqrt(3))) - (r / 2)

	// 四舍五入到最近的整数坐标
	roundedQ := int(math.Round(q))
	roundedR := int(math.Round(r))

	// 创建候选坐标
	candidate := game.Coordinate{X: roundedQ, Y: roundedR}

	// 检查该坐标是否在棋盘内
	if _, exists := game.EmptyDvonn.Cells[candidate]; exists {
		// 验证点击位置是否在这个六边形内
		candidateX, candidateY := coordToScreen(candidate)
		dx := candidateX - float64(x)
		dy := candidateY - float64(y)

		// 使用合适的点击范围判断
		if dx*dx+dy*dy <= (triangleR*0.5)*(triangleR*0.5) {
			return candidate, true
		}
	}

	return game.Coordinate{}, false
}

// 处理鼠标，返回生成的 Move（或 nil）
func handleInput(gs *game.GameState) game.Move {
	// 如果是鼠标点击，获取坐标
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		c, ok := pixelToCoord(x, y) // 将像素坐标转换为棋盘坐标
		//fmt.Println(c)              // 打印坐标调试用
		if !ok {
			return nil
		}

		// Phase1：放子阶段
		if gs.Phase == game.Phase1 {
			// 调用 runPlacementPhase 来放置棋子
			game.RunPlacementPhase(gs, c.X, c.Y)
			return nil
		}

	}
	// 无效输入或操作
	return nil
}
