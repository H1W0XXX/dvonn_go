package ai

import (
	"crypto/rand"
	"dvonn_go/internal/game"
	"encoding/binary"
)

// -------------- Zobrist 随机表 ----------------
const (
	MaxStackHeight = 49 // DVONN 初始一共 49 颗棋子，栈高不可能超过它
)

var (
	// 颜色维度：Red、White、Black
	zobristCol [3][game.BoardWidth][game.BoardHeight]uint64
	// 高度维度：1…MaxStackHeight
	zobristHgt [MaxStackHeight + 1][game.BoardWidth][game.BoardHeight]uint64
	// 逐层棋子颜色：第 idx 层 (0=顶层)
	zobristStack [MaxStackHeight][game.BoardWidth][game.BoardHeight][3]uint64
	// 走子权：表示当前轮到谁走
	zobristSide uint64
)

func init() {
	// 初始化颜色和高度维度随机表
	for p := 0; p < 3; p++ {
		for x := 0; x < game.BoardWidth; x++ {
			for y := 0; y < game.BoardHeight; y++ {
				zobristCol[p][x][y] = randomUint64()
			}
		}
	}
	for h := 1; h <= MaxStackHeight; h++ {
		for x := 0; x < game.BoardWidth; x++ {
			for y := 0; y < game.BoardHeight; y++ {
				zobristHgt[h][x][y] = randomUint64()
			}
		}
	}
	for layer := 0; layer < MaxStackHeight; layer++ {
		for x := 0; x < game.BoardWidth; x++ {
			for y := 0; y < game.BoardHeight; y++ {
				for color := 0; color < 3; color++ {
					zobristStack[layer][x][y][color] = randomUint64()
				}
			}
		}
	}
	// 初始化走子权随机数
	zobristSide = randomUint64()
}

func randomUint64() uint64 {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return binary.LittleEndian.Uint64(b[:])
}

// Hash 计算包含棋子颜色、堆叠高度和走子权的 Zobrist 哈希
func Hash(b *game.Board, turn game.Player) uint64 {
	var h uint64
	// 棋子颜色和高度
	for x := 0; x < game.BoardWidth; x++ {
		for y := 0; y < game.BoardHeight; y++ {
			//coord := game.Coordinate{X: x, Y: y} // 将坐标 (x, y) 转换为 Coordinate 类型
			st := b.Cells[x][y]
			if st == nil || len(*st) == 0 {
				continue
			}

			// 栈顶颜色索引（本项目中 stack[0] 为顶）
			top := (*st)[0] // 解引用 st，获取栈顶棋子
			topIdx := pieceIndex(top)
			height := len(*st)

			h ^= zobristCol[topIdx][x][y]
			if height <= MaxStackHeight {
				h ^= zobristHgt[height][x][y]
			}

			for layer, piece := range *st {
				if layer >= MaxStackHeight {
					break
				}
				idx := pieceIndex(piece)
				h ^= zobristStack[layer][x][y][idx]
			}
		}
	}
	// 走子权
	if turn == game.PBlack {
		h ^= zobristSide
	}
	return h
}

func pieceIndex(p game.Piece) int {
	switch p {
	case game.Red:
		return 0
	case game.White:
		return 1
	case game.Black:
		return 2
	default:
		return 0
	}
}

// ---------------- TT 结构 ---------------------

type ttEntry struct {
	depth int
	score int
	move  game.Move
}

type TT struct {
	table map[uint64]ttEntry
}

func NewTT() *TT {
	return &TT{table: make(map[uint64]ttEntry, 1<<18)}
}

func (tt *TT) Save(hash uint64, depth, score int, move game.Move) {
	e, ok := tt.table[hash]
	if !ok || depth > e.depth {
		tt.table[hash] = ttEntry{depth: depth, score: score, move: move}
	}
}

func (tt *TT) Lookup(hash uint64, depth int) (score int, move game.Move, ok bool) {
	e, ok := tt.table[hash]
	if ok && e.depth >= depth {
		return e.score, e.move, true
	}
	return 0, nil, false
}
