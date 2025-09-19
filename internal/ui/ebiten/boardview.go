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
	hexSize   = 40  // �����α߳�
	boardXOff = 300 // ����X��ƫ�ƣ�ȷ����಻���ض�
	boardYOff = 200 // ����Y��ƫ�ƣ�ȷ���Ϸ������ض�
)
const (
	layerOffsetY = 6 // ÿ�����Ӵ�ֱƫ������
)

// axial �� ��Ļ����
// coordToScreen ������ c ת��Ϊ��Ļ�ϵ�λ��
func coordToScreen(c game.Coordinate) (float64, float64) {
	q, r := axialFromIndex(c)

	// ����������ı�׼ת����ʽ
	x := triangleR * (math.Sqrt(3)*q + math.Sqrt(3)/2*r)
	y := triangleR * (3.0 / 2 * r)

	// Ӧ��ƫ������ȷ���������̶��ڿɼ�������
	return offsetX + x, offsetY + y
}

// drawStack �������Ӷѣ�α3D���+������ע
func drawStack(b *game.Board, c game.Coordinate, screen *ebiten.Image) {
	stackPtr := b.Cells[c.X][c.Y]
	if stackPtr == nil || len(*stackPtr) == 0 {
		return
	}
	stack := *stackPtr

	// ����ת��Ϊ��Ļ����
	x, y := coordToScreen(c)
	// ��ͼ���Ż�׼����������������ͼ��ͬ��С��
	scaledSize := float64(triangleR)
	scale := scaledSize / float64(imgRed.Bounds().Dx())

	// �ӵײ㵽������ƣ�ÿ����ݶ�Ӧ�� Piece ѡ����ͼ
	for idx := len(stack) - 1; idx >= 0; idx-- {
		piece := stack[idx]
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
		// ���㵱ǰ��Ĵ�ֱƫ��
		layerOffset := float64(len(stack)-1-idx) * layerOffsetY
		py := y - layerOffset
		op.GeoM.Translate(x-scaledSize/2, py-scaledSize/2)
		screen.DrawImage(img, op)
	}

	// �ڶ������Ļ��ƶѸ�����
	count := fmt.Sprint(len(stack))
	labelX := int(x) + int(scaledSize)/2 - 4
	labelY := int(y) - (len(stack)-1)*layerOffsetY + 10
	// ��ɫ��Ӱ
	text.Draw(screen, count, basicfont.Face7x13, labelX, labelY, color.Black)
	// ��ɫǰ��
	text.Draw(screen, count, basicfont.Face7x13, labelX-1, labelY-1, color.White)
}
