// File internal/ai/zobrist.go
package ai

import (
	"crypto/rand"
	"dvonn_go/internal/game"
	"encoding/binary"
)

// -------------- Zobrist 随机表 ----------------

var zobrist [3][12][6]uint64 // [Piece][X][Y]

// 初始化时生成随机 64 位
func init() {
	for p := 0; p < 3; p++ {
		for x := 0; x < 11; x++ {
			for y := 0; y < 6; y++ {
				zobrist[p][x][y] = randomUint64()
			}
		}
	}
}

func randomUint64() uint64 {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return binary.LittleEndian.Uint64(b[:])
}

// 对棋盘计算哈希
func Hash(b *game.Board) uint64 {
	var h uint64
	for c, st := range b.Cells {
		if len(st) == 0 {
			continue
		}
		// 根据堆顶颜色索引
		var idx int
		switch st[0] {
		case game.White:
			idx = 1
		case game.Black:
			idx = 2
		default:
			idx = 0
		}
		// 仅在合法索引范围内哈希
		if c.X >= 0 && c.X < 12 && c.Y >= 0 && c.Y < 6 {
			h ^= zobrist[idx][c.X][c.Y]
		}
	}
	return h
}

// ---------------- TT 结构 ---------------------

type ttEntry struct {
	depth int
	score int
	move  game.Move
}

type TT struct{ table map[uint64]ttEntry }

func NewTT() *TT { return &TT{table: make(map[uint64]ttEntry, 1<<18)} }

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
