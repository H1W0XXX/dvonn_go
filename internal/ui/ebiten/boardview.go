// File internal/ui/ebiten/boardview.go
package ebiten

import (
	"dvonn_go/internal/game"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"math"
)

const (
	hexSize   = 40  // 六边形边长
	boardXOff = 300 // 增大X轴偏移，确保左侧不被截断
	boardYOff = 200 // 增大Y轴偏移，确保上方不被截断
)
const (
	layerOffsetY = 6 // 每层棋子垂直偏移像素
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

// drawStack 绘制棋子堆，伪3D层叠+层数标注
func drawStack(b *game.Board, c game.Coordinate, screen *ebiten.Image) {
	st := b.Cells[c]
	if len(st) == 0 {
		return
	}
	// 坐标转换为屏幕像素
	x, y := coordToScreen(c)
	// 贴图缩放基准（假设所有棋子贴图相同大小）
	scaledSize := float64(triangleR)
	scale := scaledSize / float64(imgRed.Bounds().Dx())

	// 从底层到顶层绘制，每层根据对应的 Piece 选择贴图
	for idx := len(st) - 1; idx >= 0; idx-- {
		piece := st[idx]
		var img *ebiten.Image
		switch piece {
		case game.Red:
			img = imgRed
		case game.White:
			img = imgWhite
		case game.Black:
			img = imgBlack
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		// 计算当前层的垂直偏移
		layerOffset := float64(len(st)-1-idx) * layerOffsetY
		py := y - layerOffset
		op.GeoM.Translate(x-scaledSize/2, py-scaledSize/2)
		screen.DrawImage(img, op)
	}

	// 在顶层中心绘制堆高数字
	count := fmt.Sprint(len(st))
	labelX := int(x) + int(scaledSize)/2 - 4
	labelY := int(y) - (len(st)-1)*layerOffsetY + 10
	// 黑色阴影
	text.Draw(screen, count, basicfont.Face7x13, labelX, labelY, color.Black)
	// 白色前景
	text.Draw(screen, count, basicfont.Face7x13, labelX-1, labelY-1, color.White)
}
