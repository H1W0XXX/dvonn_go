// File internal/ui/ebiten/boardview.go
package ebiten

import (
	"dvonn_go/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

const (
	hexSize   = 40  // 边长
	boardXOff = 100 // 屏幕偏移
	boardYOff = 80
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
	x, y := coordToScreen(c)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y-float64(len(st))*4) // 堆高堆栈偏移
	screen.DrawImage(img, op)
}
