// File internal/ui/ebiten/gameview.go
package ebiten

import (
	"dvonn_go/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameView struct {
	state game.GameState
}

func NewGameView(gs game.GameState) *GameView {
	return &GameView{state: gs}
}

// Update 每帧处理输入与逻辑
func (g *GameView) Update() error {
	if mv := handleInput(&g.state); mv != nil {
		if game.ValidMove(&g.state.Board, mv) {
			// Phase1 顺序用 turn控制；Phase2 由 getNextTurn
			//game.ExecuteMove(&g.state, mv) // 你已在 gameflow.go 写好 executeMove
		}
	}
	return nil
}

// Draw 渲染棋盘
func (g *GameView) Draw(screen *ebiten.Image) {
	// 渲染棋盘背景
	screen.DrawImage(boardBG, nil)
	// 渲染棋子
	for c := range g.state.Board.Cells {
		drawStack(&g.state.Board, c, screen)
	}
}

func (g *GameView) Layout(outW, outH int) (int, int) { return 1300, 768 }
