// File internal/ui/ebiten/boardview.go
package ebiten

import (
	"dvonn_go/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

const (
	hexSize   = 40  // 六边形边长
	boardXOff = 300 // 增大X轴偏移，确保左侧不被截断
	boardYOff = 200 // 增大Y轴偏移，确保上方不被截断
)

// axial → 屏幕像素
// coordToScreen 将坐标 c 转换为屏幕上的位置
func coordToScreen(c game.Coordinate) (float64, float64) {
	q := float64(c.X)
	r := float64(c.Y)

	// 六边形网格的标准转换公式
	x := triangleR * (math.Sqrt(3)*q + math.Sqrt(3)/2*r)
	y := triangleR * (3.0 / 2 * r)

	// 应用偏移量，确保整个棋盘都在可见区域内
	return offsetX + x, offsetY + y
}

// DrawStack 根据堆顶颜色 / 高度绘制棋子
func drawStack(b *game.Board, c game.Coordinate, screen *ebiten.Image) {
	st := b.Cells[c]
	if len(st) == 0 {
		return
	}
	var img *ebiten.Image
	switch st[0] {
	case game.Red:
		img = imgRed
	case game.White:
		img = imgWhite
	case game.Black:
		img = imgBlack
	}

	// 计算位置
	x, y := coordToScreen(c)
	var scaledSize float64 = triangleR

	op := &ebiten.DrawImageOptions{}
	// 计算缩放比
	op.GeoM.Scale(scaledSize/float64(img.Bounds().Dx()), scaledSize/float64(img.Bounds().Dy()))

	// 动态调整垂直偏移，避免硬编码值
	// 使用 scaledSize / 2 来微调棋子居中
	op.GeoM.Translate(x-scaledSize/2, y-scaledSize/2)

	// 绘制棋子
	screen.DrawImage(img, op)
}
