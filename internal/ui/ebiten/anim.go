// File internal/ui/ebiten/anim.go
package ebiten

import (
	"dvonn_go/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

// Animation 保存一次跳子的逐帧信息
type Animation struct {
	Piece    game.Piece
	From, To game.Coordinate
	frame    int
}

// 总帧数：越大越慢
const animFrames = 10

// 完成？
func (a *Animation) done() bool { return a.frame >= animFrames }

// Draw 在屏幕上绘制当前帧的棋子
func (a *Animation) Draw(screen *ebiten.Image) {
	// 计算从 0 到 1 的 t
	t := float64(a.frame) / float64(animFrames)

	// 格点 to 屏像素
	fx, fy := coordToScreen(a.From)
	tx, ty := coordToScreen(a.To)

	// 插值坐标
	x := fx + (tx-fx)*t
	y := fy + (ty-fy)*t

	// 选图
	var img *ebiten.Image
	switch a.Piece {
	case game.Red:
		img = imgRed
	case game.White:
		img = imgWhite
	default:
		img = imgBlack
	}

	// 缩放 & 平移到中心
	op := &ebiten.DrawImageOptions{}
	scale := triangleR / float64(img.Bounds().Dx())
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x-triangleR/2, y-triangleR/2)

	// 绘制
	screen.DrawImage(img, op)
}
