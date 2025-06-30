// File internal/game/board.go
package game

import "slices"

/*
核心棋盘逻辑 —— 翻译自 Haskell 版 Board 模块
仅依赖 types.go 中的基本类型，无 UI / AI 耦合
*/

//-----------------------------------------------------------------------------
// 工具：集合操作（map[T]struct{} 充当 set）

func setAdd[S comparable](m map[S]struct{}, s S)      { m[s] = struct{}{} }
func setHas[S comparable](m map[S]struct{}, s S) bool { _, ok := m[s]; return ok }
func setDel[S comparable](m map[S]struct{}, s S)      { delete(m, s) }

//-----------------------------------------------------------------------------
// 坐标全集与合法性

// coordinates 返回棋盘上所有有效坐标
func coordinates(b *Board) map[Coordinate]struct{} {
	out := make(map[Coordinate]struct{}, len(b.Cells))
	for c := range b.Cells {
		setAdd(out, c)
	}
	return out
}

// validCoordinate 判断给定坐标是否在棋盘上
func validCoordinate(b *Board, c Coordinate) bool {
	_, ok := b.Cells[c]
	return ok
}

// allNeighbors 按 DVONN 六向邻接获取相邻坐标（若超出棋盘即丢弃）
func allNeighbors(b *Board, c Coordinate) map[Coordinate]struct{} {
	dirs := [6]Coordinate{
		{+1, 0}, {-1, 0},
		{0, +1}, {0, -1},
		{+1, +1}, {-1, -1},
	}
	ns := make(map[Coordinate]struct{}, 6)
	for _, d := range dirs {
		n := Coordinate{c.X + d.X, c.Y + d.Y}
		if validCoordinate(b, n) {
			setAdd(ns, n)
		}
	}
	return ns
}

// neighborOf: c1 是否为 c2 的邻居
func neighborOf(b *Board, c1, c2 Coordinate) bool {
	return validCoordinate(b, c2) && setHas(allNeighbors(b, c2), c1)
}

//-----------------------------------------------------------------------------
// 检查红子

// containsRed 判断单格栈内是否含红子
func containsRed(b *Board, c Coordinate) bool {
	for _, p := range b.Cells[c] {
		if p == Red {
			return true
		}
	}
	return false
}

// hasRed 判断坐标集合中是否至少含一红
func hasRed(b *Board, coords map[Coordinate]struct{}) bool {
	for c := range coords {
		if containsRed(b, c) {
			return true
		}
	}
	return false
}

//-----------------------------------------------------------------------------
// 非空/空格 & 邻居

func nonempty(b *Board, c Coordinate) bool { return len(b.Cells[c]) > 0 }

func neighbors(b *Board, c Coordinate) map[Coordinate]struct{} {
	ns := allNeighbors(b, c)
	for n := range ns {
		if !nonempty(b, n) {
			setDel(ns, n)
		}
	}
	return ns
}

//-----------------------------------------------------------------------------
// DFS 求连通块

func component(b *Board, start Coordinate) map[Coordinate]struct{} {
	if !nonempty(b, start) {
		return nil
	}
	seen := map[Coordinate]struct{}{start: {}}
	var dfs func(Coordinate)
	dfs = func(cur Coordinate) {
		for n := range neighbors(b, cur) {
			if !setHas(seen, n) {
				setAdd(seen, n)
				dfs(n)
			}
		}
	}
	dfs(start)
	return seen
}

func allComponents(b *Board) []map[Coordinate]struct{} {
	visited := make(map[Coordinate]struct{})
	var comps []map[Coordinate]struct{}
	for c := range b.Cells {
		if nonempty(b, c) && !setHas(visited, c) {
			comp := component(b, c)
			for k := range comp {
				setAdd(visited, k)
			}
			comps = append(comps, comp)
		}
	}
	return comps
}

// isSurrounded：某格若六邻全被占，则该堆不能移动
func isSurrounded(b *Board, c Coordinate) bool {
	return len(neighbors(b, c)) == 6
}

//-----------------------------------------------------------------------------
// 统计与胜负判定

// innerstack 直接取出坐标栈（nil 视为空）
func innerstack(b *Board, c Coordinate) Stack { return b.Cells[c] }

// calcWinner 根据全盘堆顶方计算胜负（高者胜，等高和局）
func calcWinner(b *Board) *Player {
	var whiteCnt, blackCnt int
	for c := range b.Cells {
		st := innerstack(b, c)
		if len(st) == 0 {
			continue
		}
		if st[0] == White {
			whiteCnt += len(st)
		} else if st[0] == Black {
			blackCnt += len(st)
		}
	}
	switch {
	case whiteCnt > blackCnt:
		p := PWhite
		return &p
	case blackCnt > whiteCnt:
		p := PBlack
		return &p
	default:
		return nil
	}
}

// nonempties / empties / countEmpty
func nonempties(b *Board) map[Coordinate]struct{} {
	s := make(map[Coordinate]struct{})
	for c := range b.Cells {
		if nonempty(b, c) {
			setAdd(s, c)
		}
	}
	return s
}

func empties(b *Board) map[Coordinate]struct{} {
	s := make(map[Coordinate]struct{})
	for c := range b.Cells {
		if !nonempty(b, c) {
			setAdd(s, c)
		}
	}
	return s
}

func countEmpty(b *Board) int { return len(empties(b)) }

// numActivePieces / numDiscardedPieces
func numActivePieces(b *Board) int {
	var n int
	for _, st := range b.Cells {
		n += len(st)
	}
	return n
}
func numDiscardedPieces(b *Board) int { return len(b.Discard) }

//-----------------------------------------------------------------------------
// 执行走子 & 清理

// place 在指定坐标顶端压入棋子
func place(b *Board, p Piece, c Coordinate) *Board {
	st := slices.Clone(b.Cells[c])
	b.Cells[c] = append([]Piece{p}, st...)
	return b
}

// combine: 将 c1 叠到 c2 顶端
func combine(b *Board, c1, c2 Coordinate) *Board {
	newStack := append(innerstack(b, c1), innerstack(b, c2)...)
	b.Cells[c1] = nil
	b.Cells[c2] = newStack
	return b
}

// discard: 把给定坐标集合中的栈扔进弃子堆，并清空原格
func discard(b *Board, coords map[Coordinate]struct{}) *Board {
	for c := range coords {
		b.Discard = append(b.Discard, b.Cells[c]...)
		b.Cells[c] = nil
	}
	return b
}

// cleanup: 移除所有与红子断开的连通块
func cleanup(b *Board) *Board {
	for _, comp := range allComponents(b) {
		if !hasRed(b, comp) {
			discard(b, comp)
		}
	}
	return b
}

// apply 根据 Move 更新棋盘；返回更新后的副本
func apply(m Move, b *Board) Board {
	switch mv := m.(type) {
	case JumpMove:
		combine(b, mv.From, mv.To) // 叠堆
		cleanup(b)                 // 断连清理
	case PlaceMove:
		place(b, mv.Piece, mv.At) // 直接放子
	}
	return *b
}
