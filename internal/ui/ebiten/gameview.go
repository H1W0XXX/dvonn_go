// File internal/ui/ebiten/gameview.go
package ebiten

import (
	"dvonn_go/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"time"
)

type GameView struct {
	state game.GameState
	mode  string
}

func NewGameView(gs game.GameState, mode string) *GameView {
	// 给随机数初始化种子
	rand.Seed(time.Now().UnixNano())
	return &GameView{state: gs, mode: mode}
}

func (g *GameView) Update() error {
	handleInput(&g.state)
	return nil
}

func (g *GameView) Draw(screen *ebiten.Image) {

	// 4. 绘制所有棋子
	for c := range g.state.Board.Cells {
		drawStack(&g.state.Board, c, screen)
	}

	// 1. 绘制棋盘背景
	screen.DrawImage(boardBG, nil)

	// 2. Phase2 且未选中时，高亮所有可移动堆（蓝色）
	if g.state.Phase == game.Phase2 && !selected {
		for _, c := range movableCoords(&g.state) {
			drawCircleColored(screen, c, highlightBlue)
		}
	}

	// 3. 已选中时，只高亮这个被选的堆（蓝色）及其所有落点（绿色）
	if selected {
		// 蓝色圈住选中的那一堆
		drawCircleColored(screen, selectedAt, highlightBlue)
		// 绿色圈出它的所有合法落点
		for _, dst := range destinations(&g.state, selectedAt) {
			drawCircleColored(screen, dst, highlightGreen)
		}
	}

}

func (g *GameView) Layout(_, _ int) (int, int) {
	return 1300, 768
}
