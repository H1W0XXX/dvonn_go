package ebiten

import (
	"dvonn_go/internal/game"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math"
)

var (
	lineColor   = color.RGBA{0x44, 0x44, 0x44, 0xFF} // 线条颜色：深灰
	circleColor = color.RGBA{0xB0, 0xC4, 0xDE, 0xFF} // 淡蓝色
	textColor   = color.RGBA{0x00, 0x00, 0x00, 0xFF}
)

// boardBG 在 init 时构建，Draw 时直接贴上
var boardBG *ebiten.Image

// 调整整体大小 / 偏移（与 boardview.go 中一致）
const (
	triangleR = 70 // 三角形边长

	offsetX   = 650
	offsetY   = 360
	fillColor = 0xDDDDDDFF // 淡灰填充

	canvasW = 1300 // 画布宽度
	canvasH = 700  // 画布高度
)

var triangleH = math.Sqrt(3) * triangleR / 2 // 等边三角形的高度
func init() {
	// 1．创建一张背景图
	boardBG = ebiten.NewImage(canvasW, canvasH)
	// 2．对每个合法格画出圆圈并连接线条
	for c := range game.EmptyDvonn.Cells {
		drawCircle(boardBG, c)
		drawCoordinate(boardBG, c)
	}
	// 绘制格子之间的线条
	drawGridLines(boardBG)
}

func drawCoordinate(dst *ebiten.Image, c game.Coordinate) {
	cx, cy := coordToScreen(c)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("(%d,%d)", c.X, c.Y), int(cx-10), int(cy-10))
}

func drawCircle(dst *ebiten.Image, c game.Coordinate) {
	screenX, screenY := coordToScreen(c)
	radius := float64(triangleR) / 3
	ebitenutil.DrawCircle(dst, screenX, screenY, radius, circleColor)
}

// 绘制六边形格子之间的连线
func drawGridLines(dst *ebiten.Image) {
	for c1 := range game.EmptyDvonn.Cells { // 遍历每个 Coordinate 键
		// 查找相邻的格子
		for _, c2 := range getAdjacentCoordinates(c1) {
			// 确保 c2 是棋盘上的有效坐标
			if _, exists := game.EmptyDvonn.Cells[c2]; exists {
				screenX1, screenY1 := coordToScreen(c1)
				screenX2, screenY2 := coordToScreen(c2)
				// 绘制连接线
				ebitenutil.DrawLine(dst, screenX1, screenY1, screenX2, screenY2, lineColor)
			}
		}
	}
}

// 获取相邻的格子坐标
func getAdjacentCoordinates(c game.Coordinate) []game.Coordinate {
	var adjacent []game.Coordinate
	// 假设相邻格子的坐标在六边形网格中相差的坐标差值
	deltas := []game.Coordinate{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, -1}, {-1, 1},
	}

	for _, delta := range deltas {
		adjacent = append(adjacent, game.Coordinate{X: c.X + delta.X, Y: c.Y + delta.Y})
	}
	return adjacent
}
