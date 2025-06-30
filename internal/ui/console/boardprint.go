// File internal/ui/console/boardprint.go
package ui

import (
	"fmt"
	"sort"

	game "dvonn_go/internal/game"
)

// -----------------------------------------------------------------------------
// 辅助：棋盘坐标与栈信息
// -----------------------------------------------------------------------------

// coordinates 返回按 (x,y) 升序排好的坐标切片
func coordinates(b *game.Board) []game.Coordinate {
	keys := make([]game.Coordinate, 0, len(b.Cells))
	for c := range b.Cells {
		keys = append(keys, c)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Y == keys[j].Y {
			return keys[i].X < keys[j].X
		}
		return keys[i].Y < keys[j].Y
	})
	return keys
}

// containsRed 判断该格是否含红子
func containsRed(b *game.Board, c game.Coordinate) bool {
	for _, p := range b.Cells[c] {
		if p == game.Red {
			return true
		}
	}
	return false
}

// innerstack 取出坐标栈
func innerstack(b *game.Board, c game.Coordinate) game.Stack { return b.Cells[c] }

// -----------------------------------------------------------------------------
// 行片段绘制
// -----------------------------------------------------------------------------

func printTops(b *game.Board, spaces, row int) {
	line := " "
	for _, c := range coordinates(b) {
		if c.Y == row {
			star := " "
			if containsRed(b, c) {
				star = "*"
			}
			line += " /" + star + "\\"
		}
	}
	fmt.Println(spacestring(spaces) + line)
}

func printLowerTops(b *game.Board, spaces, row int) {
	line := " \\"
	for _, c := range coordinates(b) {
		if c.Y == row {
			star := " "
			if containsRed(b, c) {
				star = "*"
			}
			line += " /" + star + "\\"
		}
	}
	line += "/"
	fmt.Println(spacestring(spaces) + line)
}

func printBases(b *game.Board, spaces, row int) {
	line := " "
	for _, c := range coordinates(b) {
		if c.Y == row {
			line += " \\ /"
		}
	}
	fmt.Println(spacestring(spaces) + line)
}

func printMiddlesDvonn(b *game.Board, spaces, row int) {
	line := "|"
	for _, c := range coordinates(b) {
		if c.Y == row {
			line += label(b, c) + "|"
		}
	}
	fmt.Println(spacestring(spaces)+line+spacestring(spaces)+" ", row)
}

func printTopsMini(b *game.Board, spaces, row int) {
	line := " "
	for _, c := range coordinates(b) {
		if c.Y == row {
			star := " "
			if containsRed(b, c) {
				star = "*"
			}
			line += " /" + star + "\\"
		}
	}
	line += "/"
	fmt.Println(spacestring(spaces) + line)
}

func printMiddlesMini(b *game.Board, spaces, row int) {
	line := "|"
	for _, c := range coordinates(b) {
		if c.Y == row {
			line += label(b, c) + "|"
		}
	}
	fmt.Println(spacestring(spaces)+line+spacestring(5-spaces), row)
}

// -----------------------------------------------------------------------------
// 标签 / 格内文本
// -----------------------------------------------------------------------------

func label(b *game.Board, c game.Coordinate) string {
	st := innerstack(b, c)
	if len(st) == 0 {
		return "   "
	}
	col := " "
	switch st[0] {
	case game.White:
		col = "W"
	case game.Black:
		col = "B"
	case game.Red:
		col = "R"
	}
	size := len(st)
	if size < 10 {
		return col + " " + fmt.Sprint(size)
	}
	return col + fmt.Sprint(size)
}

func printTopLabelsDvonn() {
	letters := "ABCDEFGHI"
	line := ""
	for _, ch := range letters {
		line += "   " + string(ch)
	}
	fmt.Println(spacestring(5) + line)
}

func printBottomLabelsDvonn() {
	letters := "CDEFGHIJK"
	line := ""
	for _, ch := range letters {
		line += "   " + string(ch)
	}
	fmt.Println(" " + line)
}

func printTopLabelsMini() {
	letters := "ABC"
	line := ""
	for _, ch := range letters {
		line += "   " + string(ch)
	}
	fmt.Println(spacestring(5) + line + " ")
}

func printBottomLabelsMini() { fmt.Println("A   B   C ") }

// -----------------------------------------------------------------------------
// 完整棋盘 / 迷你棋盘
// -----------------------------------------------------------------------------

func printGridDvonn(b *game.Board) {
	printTopLabelsDvonn()
	printTops(b, 4, 1)
	printMiddlesDvonn(b, 4, 1)
	printTops(b, 2, 2)
	printMiddlesDvonn(b, 2, 2)
	printTops(b, 0, 3)
	printMiddlesDvonn(b, 0, 3)
	printLowerTops(b, 0, 4)
	printMiddlesDvonn(b, 2, 4)
	printLowerTops(b, 2, 5)
	printMiddlesDvonn(b, 4, 5)
	printBases(b, 4, 5)
	printBottomLabelsDvonn()
}

func printGridMini(b *game.Board) {
	printTopLabelsMini()
	printTops(b, 4, 1)
	printMiddlesMini(b, 4, 1)
	printTopsMini(b, 2, 2)
	printMiddlesMini(b, 2, 2)
	printTopsMini(b, 0, 3)
	printMiddlesMini(b, 0, 3)
	printBases(b, 0, 3)
	printBottomLabelsMini()
}

// PrintBoard 根据棋盘大小选择合适的 ASCII 渲染
func PrintBoard(b *game.Board) {
	if len(b.Cells) == 9 {
		printGridMini(b)
	} else {
		printGridDvonn(b)
	}
}

// -----------------------------------------------------------------------------
// 小工具
// -----------------------------------------------------------------------------

func spacestring(n int) string { return string(make([]rune, n)) }
