// File internal/ui/ebiten/gameview.go
package ebiten

import (
	"dvonn_go/internal/ai"
	"dvonn_go/internal/game"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
)

const depth = 4

type GameView struct {
	state        game.GameState
	mode         string
	aiPlayer     game.Player
	anims        []*Animation
	pendingMv    *game.JumpMove //等待执行的跳子
	showedResult bool
}

func NewGameView(gs game.GameState, mode string) *GameView {
	// 约定 AI 操作黑棋，人类操作白棋
	return &GameView{
		state:        gs,
		mode:         mode,
		aiPlayer:     game.PWhite,
		anims:        []*Animation{},
		pendingMv:    nil,
		showedResult: false,
	}
}
func (g *GameView) Update() error {
	// 1) 玩家输入 → 产生 Move
	if mv := handleInput(&g.state); mv != nil {
		switch m := mv.(type) {
		case game.PlaceMove:
			// 放子阶段：立即执行
			game.RunPlacementPhase(&g.state, m.At.X, m.At.Y)

		case game.JumpMove:
			// 跳子阶段：**不** 立刻执行，只排入动画 & 挂起
			mv2 := m
			g.pendingMv = &mv2
			g.anims = append(g.anims, &Animation{
				Piece: m.Player.Piece(),
				From:  m.From,
				To:    m.To,
			})
		}
	}

	// 2) PvE AI 下子也同理，只挂起，不立刻合并
	if !game.IsGameOver(&g.state) &&
		g.mode == "pve" &&
		g.state.Phase == game.Phase2 &&
		game.TurnStateToPlayer(g.state.Turn) == g.aiPlayer &&
		clickStep == 0 &&
		len(g.anims) == 0 &&
		g.pendingMv == nil {

		best := ai.SearchBestMove(&g.state, depth)
		mv2 := best
		g.pendingMv = &mv2
		g.anims = append(g.anims, &Animation{
			Piece: best.Player.Piece(),
			From:  best.From,
			To:    best.To,
		})
	}

	// 3) 推进所有动画帧
	for _, a := range g.anims {
		a.frame++
	}

	// 4) 如果队首动画完成，再真正执行一次 RunMovementPhase
	if len(g.anims) > 0 && g.anims[0].done() && g.pendingMv != nil {
		m := *g.pendingMv
		game.RunMovementPhase(&g.state, m.From, m.To)
		g.pendingMv = nil
	}

	// 5) 清理已完成的动画
	var next []*Animation
	for _, a := range g.anims {
		if !a.done() {
			next = append(next, a)
		}
	}
	g.anims = next

	//6) 游戏结束
	if game.IsGameOver(&g.state) && !g.showedResult {
		// 1) 统计各自控制的子数
		var whiteCnt, blackCnt int
		for _, st := range g.state.Board.Cells {
			if len(st) == 0 {
				continue
			}
			switch st[0] {
			case game.White:
				whiteCnt += len(st)
			case game.Black:
				blackCnt += len(st)
			}
		}

		// 2) 输出统计结果
		fmt.Printf("Game over! White controls %d pieces; Black controls %d pieces.", whiteCnt, blackCnt)

		// 3) 输出胜负
		switch {
		case whiteCnt > blackCnt:
			fmt.Println(" White wins!")
		case blackCnt > whiteCnt:
			fmt.Println(" Black wins!")
		default:
			fmt.Println(" It's a draw!")
		}

		g.showedResult = true
	}

	return nil
}

func (g *GameView) Draw(screen *ebiten.Image) {
	// 1. 绘制棋盘背景
	screen.DrawImage(boardBG, nil)

	// 2. Phase2 且未选中时，高亮可移动堆
	if g.state.Phase == game.Phase2 && !selected {
		for _, c := range movableCoords(&g.state) {
			drawCircleColored(screen, c, highlightBlue)
		}
	}

	// 3. 已选中时，高亮起点及落点
	if selected {
		drawCircleColored(screen, selectedAt, highlightBlue)
		for _, dst := range destinations(&g.state, selectedAt) {
			drawCircleColored(screen, dst, highlightGreen)
		}
	}

	// 4. 计算 “隐藏” 列表 ——  在动画播放期间，隐藏源格和目标格的静态绘制
	hide := map[game.Coordinate]bool{}
	for _, a := range g.anims {
		hide[a.From] = true
	}

	// 5. 绘制所有静态棋子（跳过正在动画的源/目标格）
	for c := range g.state.Board.Cells {
		if !hide[c] {
			drawStack(&g.state.Board, c, screen)
		}
	}

	// 6. 播放动画帧
	for _, a := range g.anims {
		a.Draw(screen)
	}
}
func (g *GameView) Layout(_, _ int) (int, int) {
	return 1300, 768
}
