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

// handleInput 仅负责返回用户的一步操作：
//
//	Phase1 → PlaceMove{Piece, At}
//	Phase2 → JumpMove{Player, From, To}
//
// 未按出合法一步时返回 nil。
// handleInput 只负责从鼠标点击生成一个 Move；
// Phase1 → PlaceMove，Phase2 → JumpMove；
// 不再在此直接调用游戏逻辑或动画。
func handleInput(gs *game.GameState) game.Move {
	// 只在“刚按下”触发
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return nil
	}

	x, y := ebiten.CursorPosition()
	c, ok := pixelToCoord(x, y)
	if !ok {
		// 点击在棋盘外，取消选中
		clickStep = 0
		selected = false
		return nil
	}

	// Phase 1：放子
	if gs.Phase == game.Phase1 {
		game.RunPlacementPhase(gs, c.X, c.Y)
		return nil
	}

	// Phase2：跳子
	if gs.Phase == game.Phase2 {
		switch clickStep {
		case 0:
			// 第一次按：选起点
			if !isMovable(gs, c) {
				return nil
			}
			fromCoord = c
			clickStep = 1
			selected = true
			selectedAt = c
			return nil

		case 1:
			// 第二次按：选终点
			for _, d := range destinations(gs, fromCoord) {
				if d == c {
					enterPerf()
					// 构造并返回 JumpMove
					pl := game.TurnStateToPlayer(gs.Turn)
					clickStep = 0
					selected = false
					return game.JumpMove{Player: pl, From: fromCoord, To: c}
				}
			}
			// 非法落点，取消选中
			clickStep = 0
			selected = false
			return nil
		}
	}

	return nil
}
