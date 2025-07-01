// File internal/game/board.go
package game

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
		{+1, 0},
		{+1, -1},
		{0, -1},
		{-1, 0},
		{-1, +1},
		{0, +1},
	}
	ns := make(map[Coordinate]struct{}, 6)
	for _, d := range dirs {
		n := Coordinate{c.X + d.X, c.Y + d.Y}
		if validCoordinate(b, n) {
			ns[n] = struct{}{}
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

func neighbors(b *Board, c Coordinate) []Coordinate {
	// 保护：未初始化就直接返回空
	if b == nil || b.Cells == nil {
		return nil
	}

	// 六个轴向相邻方向
	dirs := [6]Coordinate{
		{+1, 0},
		{+1, -1},
		{0, -1},
		{-1, 0},
		{-1, +1},
		{0, +1},
	}

	var result []Coordinate
	for _, d := range dirs {
		nc := Coordinate{X: c.X + d.X, Y: c.Y + d.Y}
		// 只把真正存在于 b.Cells 中的格子当邻居
		if _, ok := b.Cells[nc]; ok {
			result = append(result, nc)
		}
	}
	return result
}

// -----------------------------------------------------------------------------
// DFS 求连通块
func component(b *Board, start Coordinate) map[Coordinate]struct{} {
	if !nonempty(b, start) {
		return nil
	}
	seen := map[Coordinate]struct{}{start: {}}
	var dfs func(Coordinate)
	dfs = func(cur Coordinate) {
		for _, n := range neighbors(b, cur) {
			// 只对有棋子的邻居继续遍历
			if nonempty(b, n) {
				if _, visited := seen[n]; !visited {
					seen[n] = struct{}{}
					dfs(n)
				}
			}
		}
	}
	dfs(start)
	return seen
}

// allComponents 扫描所有非空格子，对每个还没访问过的格子用 component 拆出一个连通块。
func allComponents(b *Board) []map[Coordinate]struct{} {
	visited := make(map[Coordinate]struct{})
	var comps []map[Coordinate]struct{}

	for c := range b.Cells {
		if nonempty(b, c) {
			// 跳过已经归入某个连通块的格子
			if _, seen := visited[c]; !seen {
				comp := component(b, c)
				// 标记所有 comp 中的格子已访问
				for k := range comp {
					visited[k] = struct{}{}
				}
				comps = append(comps, comp)
			}
		}
	}
	return comps
}

// isSurrounded：某格若六邻全被占，则该堆不能移动
func isSurrounded(b *Board, c Coordinate) bool {
	if b == nil || b.Cells == nil {
		return false
	}

	// 六个轴向相邻方向向量
	dirs := []Coordinate{
		{+1, -1}, {+1, 0}, {0, +1},
		{-1, +1}, {-1, 0}, {0, -1},
	}

	for _, d := range dirs {
		nc := Coordinate{X: c.X + d.X, Y: c.Y + d.Y}
		stack, exists := b.Cells[nc]
		// 只要有一个方向：①不在棋盘（exists==false）
		// 或者 ②在棋盘但空（len(stack)==0），就不是被包围
		if !exists || len(stack) == 0 {
			return false
		}
	}
	return true
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
// place 在指定坐标顶端压入棋子；如坐标非法或已被占则保持原状
func place(b *Board, p Piece, c Coordinate) *Board {
	// 坐标非法
	if !validCoordinate(b, c) {
		return b
	}
	// 该格已被占
	if nonempty(b, c) {
		return b
	}

	// 压栈
	b.Cells[c] = []Piece{p}
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

// -----------------------------------------------------------------------------
// Apply —— 将 Move 应用到棋盘，返回更新后的 *值拷贝*
// -----------------------------------------------------------------------------
// 约定：调用者已用 ValidMove 判断合法性；这里不再做额外校验。
func Apply(m Move, b *Board) Board {
	switch mv := m.(type) {
	case JumpMove:
		// ① 把 mv.From 叠到 mv.To
		combine(b, mv.From, mv.To)
		// ② 移除与红子断开的连通块
		cleanup(b)

	case PlaceMove:
		// 直接在目标顶端放入棋子
		place(b, mv.Piece, mv.At)

	default:
		// 未知 Move 类型 —— 不做任何改动
	}

	// 返回修改后的「值拷贝」，方便写：gs.Board = Apply(mv, &gs.Board)
	return *b
}

// Clone 深度拷贝一个 Board，Cells map 与 Discard 切片都新建
func (b *Board) Clone() Board {
	// 复制 Cells
	newCells := make(map[Coordinate]Stack, len(b.Cells))
	for c, st := range b.Cells {
		// st 是 []Piece
		newSt := make([]Piece, len(st))
		copy(newSt, st)
		newCells[c] = newSt
	}
	// 复制 Discard
	newDiscard := make([]Piece, len(b.Discard))
	copy(newDiscard, b.Discard)

	return Board{
		Cells:   newCells,
		Discard: newDiscard,
	}
}
