// File cmd/dvonn-gui/main.go
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"dvonn_go/internal/game"         // 规则与状态
	ui "dvonn_go/internal/ui/ebiten" // GUI 渲染 / 输入层
)

func main() {
	// 创建初始 GameView（持有 GameState）
	view := ui.NewGameView(game.StartState())

	// 窗口参数与 GameView.Layout 对齐
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("DVONN – Ebiten GUI")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	if err := ebiten.RunGame(view); err != nil {
		log.Fatal(err)
	}
}
