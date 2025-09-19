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

// validCoordinate 判断给定坐标是否在棋盘上
func validCoordinate(b *Board, c Coordinate) bool {
	// 保证坐标在定长数组的范围内
	return c.X >= 0 && c.X < BoardWidth && c.Y >= 0 && c.Y < BoardHeight
}

// coordinates 返回棋盘上所有有效坐标

func coordinates(b *Board) map[Coordinate]struct{} {
	out := make(map[Coordinate]struct{})
	for x := 0; x < BoardWidth; x++ {
		for y := 0; y < BoardHeight; y++ {
			setAdd(out, Coordinate{X: x, Y: y})
		}
	}
	return out
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
	ns := make(map[Coordinate]struct{})
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
	stack := b.Cells[c.X][c.Y]
	if stack == nil {
		return false
	}
	for _, p := range *stack {
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

func nonempty(b *Board, c Coordinate) bool {
	stack := b.Cells[c.X][c.Y]
	return stack != nil && len(*stack) > 0
}

func neighbors(b *Board, c Coordinate) []Coordinate {
	if b == nil {
		return nil
	}

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
		if validCoordinate(b, nc) {
			result = append(result, nc)
		}
	}
	return result
}

// -----------------------------------------------------------------------------
// allComponents 扫描所有非空格子，对每个还没访问过的格子用 component 拆出一个连通块。
func component(b *Board, start Coordinate) map[Coordinate]struct{} {
	if !nonempty(b, start) {
		return nil
	}
	seen := map[Coordinate]struct{}{start: {}}
	var dfs func(Coordinate)
	dfs = func(cur Coordinate) {
		for _, n := range neighbors(b, cur) {
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

// isSurrounded：某格若六邻全被占，则该堆不能移动
func isSurrounded(b *Board, c Coordinate) bool {
	dirs := []Coordinate{
		{+1, -1}, {+1, 0}, {0, +1},
		{-1, +1}, {-1, 0}, {0, -1},
	}

	for _, d := range dirs {
		nc := Coordinate{X: c.X + d.X, Y: c.Y + d.Y}
		if !validCoordinate(b, nc) {
			return false
		}
		stack := b.Cells[nc.X][nc.Y]
		if stack == nil || len(*stack) == 0 {
			return false
		}
	}
	return true
}

//-----------------------------------------------------------------------------
// 统计与胜负判定

// innerstack 直接取出坐标栈（nil 视为空）
func innerstack(b *Board, c Coordinate) Stack {
	stack := b.Cells[c.X][c.Y]
	if stack == nil {
		return nil
	}
	return *stack
}

// calcWinner 根据全盘堆顶方计算胜负（高者胜，等高和局）
func calcWinner(b *Board) *Player {
	var whiteCnt, blackCnt int
	for x := 0; x < BoardWidth; x++ {
		for y := 0; y < BoardHeight; y++ {
			c := Coordinate{X: x, Y: y}
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
	for x := 0; x < BoardWidth; x++ {
		for y := 0; y < BoardHeight; y++ {
			c := Coordinate{X: x, Y: y}
			if nonempty(b, c) {
				setAdd(s, c)
			}
		}
	}
	return s
}

func empties(b *Board) map[Coordinate]struct{} {
	s := make(map[Coordinate]struct{})
	for x := 0; x < BoardWidth; x++ {
		for y := 0; y < BoardHeight; y++ {
			c := Coordinate{X: x, Y: y}
			if !nonempty(b, c) {
				setAdd(s, c)
			}
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

// place 在指定坐标顶端压入棋子；如坐标非法或已被占则保持原状
func place(b *Board, p Piece, c Coordinate) *Board {
	if !validCoordinate(b, c) {
		return b
	}

	// 如果该位置已有棋子，则不做修改
	if nonempty(b, c) {
		return b
	}

	// 放置棋子
	b.Cells[c.X][c.Y] = &Stack{p}
	return b
}

// combine: 将 c1 叠到 c2 顶端
func combine(b *Board, c1, c2 Coordinate) *Board {
	stack1 := b.Cells[c1.X][c1.Y]
	stack2 := b.Cells[c2.X][c2.Y]
	if stack1 == nil || stack2 == nil {
		return b
	}

	newStack := append(*stack1, *stack2...)
	b.Cells[c2.X][c2.Y] = &newStack
	b.Cells[c1.X][c1.Y] = nil
	return b
}

// discard: 把给定坐标集合中的栈扔进弃子堆，并清空原格
func discard(b *Board, coords map[Coordinate]struct{}) *Board {
	for c := range coords {
		stack := b.Cells[c.X][c.Y]
		if stack != nil {
			b.Discard = append(b.Discard, *stack...)
			b.Cells[c.X][c.Y] = nil
		}
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
func allComponents(b *Board) []map[Coordinate]struct{} {
	visited := make(map[Coordinate]struct{}) // 用来记录已访问过的格子
	var comps []map[Coordinate]struct{}      // 存储所有的连通块

	for x := 0; x < BoardWidth; x++ {
		for y := 0; y < BoardHeight; y++ {
			c := Coordinate{X: x, Y: y}
			if nonempty(b, c) {
				// 如果该格子非空且尚未访问过
				if _, seen := visited[c]; !seen {
					comp := component(b, c) // 获取该连通块
					// 标记所有 comp 中的格子已访问
					for k := range comp {
						visited[k] = struct{}{}
					}
					comps = append(comps, comp) // 将连通块添加到结果中
				}
			}
		}
	}

	return comps
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
	// 使用定长数组来创建 newCells
	var newCells [BoardWidth][BoardHeight]*Stack

	for x := 0; x < BoardWidth; x++ {
		for y := 0; y < BoardHeight; y++ {
			stack := b.Cells[x][y]
			if stack != nil {
				// 复制 Stack
				newStack := append([]Piece(nil), *stack...)
				newCells[x][y] = (*Stack)(&newStack)
			}
		}
	}

	return Board{
		Cells:   newCells,
		Discard: append([]Piece(nil), b.Discard...), // 复制 Discard
	}
}

// 返回所有初始的红子坐标
func (b *Board) GetSourceCoordinates() []Coordinate {
	var src []Coordinate
	for x := 0; x < BoardWidth; x++ {
		for y := 0; y < BoardHeight; y++ {
			st := b.Cells[x][y]
			if st != nil {
				for _, p := range *st { // 解引用 st，遍历栈中的每个 Piece
					if p == Red { // 比较 Piece 是否是 Red
						src = append(src, Coordinate{X: x, Y: y})
						break
					}
				}
			}
		}
	}
	return src
}

// 计算棋盘上两点最远距离（比如枚举两两 HexDistance 取最大）
func (b *Board) BoardDiameter() int {
	max := 0
	// 简单 O(n²) 也行，棋盘格数很少
	var coords []Coordinate
	for x := 0; x < BoardWidth; x++ {
		for y := 0; y < BoardHeight; y++ {
			coords = append(coords, Coordinate{X: x, Y: y})
		}
	}

	for i := 0; i < len(coords); i++ {
		for j := i + 1; j < len(coords); j++ {
			d := HexDistance(coords[i], coords[j])
			if d > max {
				max = d
			}
		}
	}

	return max
}

// HexDistance 计算两个六边形格点之间的距离（六边形格距离）
// 坐标系：使用轴坐标 (q=a.X, r=a.Y)
// 转为立方坐标 (x=q, z=r, y=-x-z)，再取三个轴向差值的最大值
func HexDistance(a, b Coordinate) int {
	// axial -> cube
	x1, z1 := a.X, a.Y
	y1 := -x1 - z1
	x2, z2 := b.X, b.Y
	y2 := -x2 - z2

	// 差值
	dx := abs(x1 - x2)
	dy := abs(y1 - y2)
	dz := abs(z1 - z2)

	// 最大值即为格距
	if dx > dy && dx > dz {
		return dx
	}
	if dy > dz {
		return dy
	}
	return dz
}
