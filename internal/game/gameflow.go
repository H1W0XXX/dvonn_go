// File internal/game/gameflow.go
package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

var order = []Piece{Red, Red, Red}

const totalPieceNum = 49

func init() {

	// 交替放置黑白棋子
	l := len(order)
	for i := 0; i < totalPieceNum-l; i++ {
		if i%2 == 0 {
			order = append(order, Black)
		} else {
			order = append(order, White)
		}
	}
}

// -----------------------------------------------------------------------------
// 初始状态 & 基础工具
// -----------------------------------------------------------------------------
// StartState 返回一局新 DVONN 的初始状态
func StartState() GameState {
	return GameState{
		Board:     EmptyDvonn,
		Turn:      PlacingRed,
		Phase:     Phase1,
		PlaceStep: 0,
	}
}

func prompt(t TurnState) string {
	switch t {
	case PlacingRed:
		return "Place a red piece"
	case PlacingWhite:
		return "Place a white piece"
	case PlacingBlack:
		return "Place a black piece"
	case MoveWhite:
		return "White move"
	case MoveBlack:
		return "Black move"
	default:
		return "UnknownTurnState"
	}
}

func readLine() string {
	r := bufio.NewReader(os.Stdin)
	s, _ := r.ReadString('\n')
	return strings.TrimSpace(s)
}

func printErr(e GameError) {
	switch e {
	case MoveParseError:
		fmt.Println("You entered a malformed move.")
	case InvalidMove:
		fmt.Println("Your move was incorrect.")
	default:
		fmt.Println("Unknown error.")
	}
}
func getTurnInput(gs *GameState) Move {
	// We don't use text input in the GUI version; all moves
	// are generated by mouse events in the Ebiten layer.
	return nil
}

// -----------------------------------------------------------------------------
// 执行走子
// -----------------------------------------------------------------------------
//func ExecuteMove(gs *GameState, mv Move) {
//	// 验证走子是否合法
//	if !ValidMove(&gs.Board, mv) {
//		printErr(InvalidMove)
//		return
//	}
//
//	// 应用走子，返回新的棋盘
//	gs.Board = Apply(mv, &gs.Board)
//
//	switch m := mv.(type) {
//	case JumpMove:
//		// 更新下一回合
//		gs.Turn = GetNextTurn(&gs.Board, m.Player)
//	case PlaceMove:
//		// 放子阶段：由 runPlacementPhase 控制顺序，不在这里更新
//	}
//
//	// 检查是否有合法的跳子，如果没有合法跳子则结束回合
//	if !HasAnyLegalMoves(&gs.Board, gs.Turn) {
//		gs.Turn = End
//	}
//}

// -----------------------------------------------------------------------------
// Phase 1 — 放子
// -----------------------------------------------------------------------------
func RunPlacementPhase(gs *GameState, x, y int) {
	// 已经放满 49 子就直接返回
	if gs.PlaceStep >= totalPieceNum {
		return
	}

	piece := order[gs.PlaceStep] // 取本步应放棋色
	coord := Coordinate{X: x, Y: y}

	// 若该格已被占或坐标非法则忽略点击
	if !validCoordinate(&gs.Board, coord) || nonempty(&gs.Board, coord) {
		return
	}

	// 真正放子
	place(&gs.Board, piece, coord)

	// 步数 +1
	gs.PlaceStep++

	// 放满 49 子后进入 Phase2，按规则黑方先行
	if gs.PlaceStep >= totalPieceNum {
		gs.Phase = Phase2
		gs.Turn = MoveBlack
	}
}

// FillPhase1Auto 自动填充第一阶段棋子
func FillPhase1Auto(gs *GameState) GameState {
	// 先自动填充第一阶段的棋子，模拟简单的 AI 自动放子
	for gs.PlaceStep < totalPieceNum {
		// 假设自动填充位置的逻辑（随机选择或特定策略）
		x, y := findEmptySpot(&gs.Board)
		RunPlacementPhase(gs, x, y)
	}
	return *gs
}

// findEmptySpot 找到棋盘上一个空格
func findEmptySpot(b *Board) (int, int) {
	// 获取所有空的坐标
	var emptyCoords []Coordinate
	for x := -5; x <= 5; x++ {
		for y := -2; y <= 2; y++ {
			coord := Coordinate{X: x, Y: y}
			if !nonempty(b, coord) { // 如果该位置为空
				emptyCoords = append(emptyCoords, coord)
			}
		}
	}

	// 如果有空格，随机选择一个空格
	if len(emptyCoords) > 0 {
		randomIndex := rand.Intn(len(emptyCoords)) // 从空格中随机选一个
		return emptyCoords[randomIndex].X, emptyCoords[randomIndex].Y
	}

	// 如果没有空格，返回一个默认坐标
	return -5, -2 // 或者根据实际情况返回其他默认值
}

// -----------------------------------------------------------------------------
// Phase 2 — 跳子
// -----------------------------------------------------------------------------
func RunMovementPhase(gs *GameState, from, to Coordinate) {
	// 当前行动方
	pl := TurnStateToPlayer(gs.Turn)

	// 构造跳子动作
	mv := JumpMove{Player: pl, From: from, To: to}

	// 合法性校验
	if !ValidMove(&gs.Board, mv) {
		fmt.Printf("Invalid move: player=%v, from stack=%v\n", mv.Player, gs.Board.Cells[mv.From])
		return
	}

	// 应用跳子（combine+cleanup 都在 Apply 里完成）
	gs.Board = Apply(mv, &gs.Board)

	// 更新下一手
	gs.Turn = GetNextTurn(&gs.Board, pl)

	// 检查是否还有合法跳子
	if !HasAnyLegalMoves(&gs.Board, gs.Turn) {
		gs.Turn = End
	}
}
func TurnStateToPlayer(ts TurnState) Player {
	switch ts {
	case MoveWhite, PlacingWhite:
		return PWhite
	case MoveBlack, PlacingBlack:
		return PBlack
	default:
		// 放子阶段红棋或 Start/End 等——默认给 White，反正不会用到
		return PWhite
	}
}

// IsGameOver 判断当前局面是否结束
func IsGameOver(gs *GameState) bool {
	return gs.Turn == End
}

// Winner 返回赢家指针；平局返回 nil。
// 仅在 IsGameOver==true 时调用。
func Winner(gs *GameState) *Player {
	return calcWinner(&gs.Board) // 已在 board.go 实现
}

var summary1 = `Welcome to DVONN.

In Phase 1, you must place your pieces one by one on
unoccupied spaces on the game board. White will begin
and players will alternate placing. Place your red
pieces first, and then proceed to place your normal pieces.
White has two red pieces and Black has one red piece.

To place a piece, use the command <coord>, for example, "A3".`

var summary2 = `
In Phase 2, you may move any stack of height n with your color
piece on top to be on top of another stack that is exactly
n hops away. Be careful, though! Any stacks that are not
in a connected component with a red piece will be discarded.
Whoever controls the most pieces at the end wins!

To jump, use the command <start> "to" <end>, for example "A1 to A3".`
