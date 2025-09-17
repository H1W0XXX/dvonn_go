package ai

import (
	"crypto/rand"
	"dvonn_go/internal/game"
	"encoding/binary"
)

// -------------- Zobrist 随机表 ----------------
const (
	BoardW, BoardH = 12, 6
	MaxStackHeight = 49 // DVONN 初始一共 49 颗棋子，栈高不可能超过它
)

var (
	// 颜色维度：Red、White、Black
	zobristCol [3][BoardW][BoardH]uint64
	// 高度维度：1…MaxStackHeight
	zobristHgt [MaxStackHeight + 1][BoardW][BoardH]uint64
	// 走子权：表示当前轮到谁走
	zobristSide uint64
)

func init() {
	// 初始化颜色和高度维度随机表
	for p := 0; p < 3; p++ {
		for x := 0; x < BoardW; x++ {
			for y := 0; y < BoardH; y++ {
				zobristCol[p][x][y] = randomUint64()
			}
		}
	}
	for h := 1; h <= MaxStackHeight; h++ {
		for x := 0; x < BoardW; x++ {
			for y := 0; y < BoardH; y++ {
				zobristHgt[h][x][y] = randomUint64()
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
	for c, st := range b.Cells {
		if len(st) == 0 {
			continue
		}
		// 栈顶颜色索引（本项目中 stack[0] 为顶）
		top := st[0]
		var colorIdx int
		switch top {
		case game.Red:
			colorIdx = 0
		case game.White:
			colorIdx = 1
		case game.Black:
			colorIdx = 2
		}
		// 栈高
		height := len(st)
		if c.X >= 0 && c.X < BoardW && c.Y >= 0 && c.Y < BoardH {
			h ^= zobristCol[colorIdx][c.X][c.Y]
			if height <= MaxStackHeight {
				h ^= zobristHgt[height][c.X][c.Y]
			}
		}
	}
	// 走子权
	if turn == game.PBlack {
		h ^= zobristSide
	}
	return h
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
