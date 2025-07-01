// File internal/ui/ebiten/input.go
package ebiten

import (
	"dvonn_go/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
)

type dragState struct {
	active bool
	from   game.Coordinate
}

var d dragState
var (
	// clickStep = 0   还未选起点
	// clickStep = 1   已选起点，等待终点
	clickStep int

	// 记录 Phase2 用户第一次点击的坐标
	fromCoord game.Coordinate
)
var (
	selected   bool            // 是否已选定起点
	selectedAt game.Coordinate // 记录选中的堆
)

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

// handleInput 统一处理鼠标事件
func handleInput(gs *game.GameState) {
	// 只在鼠标“刚按下”时触发
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}

	x, y := ebiten.CursorPosition()
	c, ok := pixelToCoord(x, y)
	if !ok {
		// 点击在棋盘外，取消选中
		selected = false
		clickStep = 0
		return
	}

	// Phase 1：放子
	if gs.Phase == game.Phase1 {
		game.RunPlacementPhase(gs, c.X, c.Y)
		return
	}

	// Phase 2：跳子
	if gs.Phase == game.Phase2 {
		// 第一次点击：尝试选起点
		if clickStep == 0 {
			if !isMovable(gs, c) {
				return // 非可动堆无效
			}
			fromCoord = c
			clickStep = 1
			selected = true
			selectedAt = c
			return
		}

		// 第二次点击：如果是合法落点则跳子，否则取消选中
		dests := destinations(gs, fromCoord)
		valid := false
		for _, d := range dests {
			if d == c {
				valid = true
				break
			}
		}
		if valid {
			game.RunMovementPhase(gs, fromCoord, c)
		}
		// 无论是跳子成功还是点击非落点，都重置选中状态
		selected = false
		clickStep = 0
	}
}
