// File: internal/ui/ebiten/gameview.go
package ebiten

import (
	"dvonn_go/internal/ai"
	"dvonn_go/internal/game"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"image/color"
)

const depth = 4

type GameView struct {
	state        game.GameState
	mode         string
	aiPlayer     game.Player
	anims        []*Animation
	pendingMv    *game.JumpMove //�ȴ�ִ�е�����
	showedResult bool

	// ��������ǡ���ǰ���׶����Ƿ�����AI���������ڶ������Ž����󴥷�ʡ��
	aiAnimPlaying bool
}

func NewGameView(gs game.GameState, mode string) *GameView {
	// Լ�� AI �������壬�����������
	return &GameView{
		state:         gs,
		mode:          mode,
		aiPlayer:      game.PWhite,
		anims:         []*Animation{},
		pendingMv:     nil,
		showedResult:  false,
		aiAnimPlaying: false,
	}
}

func (g *GameView) Update() error {
	// 1) ������� �� ���� Move
	if mv := handleInput(&g.state); mv != nil {
		switch m := mv.(type) {
		case game.PlaceMove:
			// ���ӽ׶Σ�����ִ��
			game.RunPlacementPhase(&g.state, m.At.X, m.At.Y)

		case game.JumpMove:
			// ���ӽ׶Σ�**��** ����ִ�У�ֻ���붯�� & ������Ҷ����������� aiAnimPlaying��
			mv2 := m
			g.pendingMv = &mv2
			g.anims = append(g.anims, &Animation{
				Piece: m.Player.Piece(),
				From:  m.From,
				To:    m.To,
			})
		}
	}

	// 2) PvE AI ���ӣ�ֻ���𣬲����̺ϲ�
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
		// ��ǣ���ǰ���ŵ���һ�ζ������� AI
		g.aiAnimPlaying = true
	}

	// 3) �ƽ����ж���֡
	for _, a := range g.anims {
		a.frame++
	}

	// 4) ������׶�����ɣ�������ִ��һ�� RunMovementPhase
	if len(g.anims) > 0 && g.anims[0].done() && g.pendingMv != nil {
		m := *g.pendingMv
		game.RunMovementPhase(&g.state, m.From, m.To)
		g.pendingMv = nil
	}

	// 5) ��������ɵĶ���
	//    �����ȷʵ�Ӷ������Ƴ�������ɶ������Ҹöζ������� AI��
	//    ����**�������Ž�����**���� leavePerf()
	{
		had := len(g.anims)
		var next []*Animation
		for _, a := range g.anims {
			if !a.done() {
				next = append(next, a)
			}
		}
		removed := had > len(next) // ��֡�ж������������Ƴ�
		g.anims = next

		if removed && g.aiAnimPlaying {
			// ���� AI ���Ƕζ���������Ϻ󴥷�һ��ʡ��
			g.aiAnimPlaying = false
			leavePerf()
		}
	}

	//6) ��Ϸ����
	if game.IsGameOver(&g.state) && !g.showedResult {
		// 1) ͳ�Ƹ��Կ��Ƶ�����
		var whiteCnt, blackCnt int
		forEachCoordinate(func(c game.Coordinate) {
			st := g.state.Board.Cells[c.X][c.Y]
			if st == nil || len(*st) == 0 {
				return
			}
			switch (*st)[0] {
			case game.White:
				whiteCnt += len(*st)
			case game.Black:
				blackCnt += len(*st)
			}
		})

		// 2) ���ͳ�ƽ��
		fmt.Printf("Game over! White controls %d pieces; Black controls %d pieces.", whiteCnt, blackCnt)

		// 3) ���ʤ��
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

	booted = true
	// ?? ȥ��ԭ��ÿ֡���������� leavePerf() �Ĵ���
	return nil
}

func (g *GameView) Draw(screen *ebiten.Image) {
	// 1. �������̱���
	screen.DrawImage(boardBG, nil)

	// 1.1 ���Ͻ���ʾ˫����ǰ�ɿ�������
	blackScore, whiteScore := currentScores(&g.state.Board)
	drawScoreboard(screen, blackScore, whiteScore)

	// 2. Phase2 ��δѡ��ʱ���������ƶ���
	if g.state.Phase == game.Phase2 && !selected {
		for _, c := range movableCoords(&g.state) {
			drawCircleColored(screen, c, highlightBlue)
		}
	}

	// 3. ��ѡ��ʱ��������㼰���
	if selected {
		drawCircleColored(screen, selectedAt, highlightBlue)
		for _, dst := range destinations(&g.state, selectedAt) {
			drawCircleColored(screen, dst, highlightGreen)
		}
	}

	// 4. ���� �����ء� �б� ����  �ڶ��������ڼ䣬����Դ���Ŀ���ľ�̬����
	hide := map[game.Coordinate]bool{}
	for _, a := range g.anims {
		hide[a.From] = true
	}

	// 5. �������о�̬���ӣ��������ڶ�����Դ/Ŀ���
	forEachCoordinate(func(c game.Coordinate) {
		if hide[c] {
			return
		}
		drawStack(&g.state.Board, c, screen)
	})

	// 6. ���Ŷ���֡
	for _, a := range g.anims {
		a.Draw(screen)
	}
}
func (g *GameView) Layout(_, _ int) (int, int) {
	return 1300, 768
}

// currentScores ͳ�ƺڰ�˫����ǰ���Ƶ���������
func currentScores(b *game.Board) (black, white int) {
	forEachCoordinate(func(c game.Coordinate) {
		st := b.Cells[c.X][c.Y]
		if st == nil || len(*st) == 0 {
			return
		}
		switch (*st)[0] {
		case game.Black:
			black += len(*st)
		case game.White:
			white += len(*st)
		}
	})
	return
}

func drawScoreboard(screen *ebiten.Image, blackScore, whiteScore int) {
	const (
		marginX     = 20
		marginY     = 30
		lineSpacing = 20
	)
	shadowColor := color.Black
	textColor := color.White

	labelBlack := fmt.Sprintf("Black: %d", blackScore)
	labelWhite := fmt.Sprintf("White: %d", whiteScore)

	drawTextWithShadow(screen, labelBlack, marginX, marginY, shadowColor, textColor)
	drawTextWithShadow(screen, labelWhite, marginX, marginY+lineSpacing, shadowColor, textColor)
}

func drawTextWithShadow(screen *ebiten.Image, label string, x, y int, shadowColor, textColor color.Color) {
	text.Draw(screen, label, basicfont.Face7x13, x+1, y+1, shadowColor)
	text.Draw(screen, label, basicfont.Face7x13, x, y, textColor)
}
