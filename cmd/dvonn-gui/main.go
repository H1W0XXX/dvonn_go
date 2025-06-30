// File cmd/dvonn-gui/main.go
package main

import (
	"flag"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"dvonn_go/internal/game"         // 规则与状态
	ui "dvonn_go/internal/ui/ebiten" // GUI 渲染 / 输入层
)

var autoPlace bool
var mode string

func init() {
	// 解析命令行参数
	flag.BoolVar(&autoPlace, "auto", false, "是否自动填充第一阶段棋子 (default: false)")
	flag.StringVar(&mode, "mode", "pve", "游戏模式：pvp 或 pve (default: pvp)")
	flag.Parse()
}

func main() {
	// 创建初始状态，传入模式控制自动填充棋子与模式选择
	gs := game.StartState()

	// 如果是 PvE 模式并且启用了自动填充，则自动填充第一阶段棋子
	if mode == "pve" {
		gs = game.FillPhase1Auto(&gs)
	}

	// 如果是 PvP 模式且启用了自动填充，则由游戏逻辑自动放置第一阶段的棋子
	if mode == "pvp" && autoPlace {
		gs = game.FillPhase1Auto(&gs)
	}

	// 创建初始 GameView（持有 GameState）
	view := ui.NewGameView(gs)

	// 窗口参数与 GameView.Layout 对齐
	ebiten.SetWindowSize(1024, 500)
	ebiten.SetWindowTitle("DVONN – Ebiten GUI")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	// 运行游戏
	if err := ebiten.RunGame(view); err != nil {
		ebitenutil.DebugPrint(nil, err.Error())
		log.Fatal(err)
	}
}
