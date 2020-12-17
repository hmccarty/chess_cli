package main

import (
	//"fmt"
)

// Enum to define piece classes (or types)
type PieceType int
const (
	KING PieceType = iota
	QUEEN
	ROOK
	BISHOP
	KNIGHT
	PAWN
)

type Game struct {
	whiteBoard [6]uint64
	blackBoard [6]uint64
}

func (game *Game) Setup() {
	// Add pieces to bitboards, bitshifting math is
	// redundant but better represents positions within
	// the bitmap

	// Add kings on 5th file
	game.whiteBoard[KING] = 0x08
	game.blackBoard[KING] = 0x08 << 56

	// Add queens on 4th file
	game.whiteBoard[QUEEN] = 0x10
	game.blackBoard[QUEEN] = 0x10 << 56

	// Add rooks on 1st and 8th file
	game.whiteBoard[ROOK] = 0x81
	game.blackBoard[ROOK] = 0x81 << 56

	// Add bishops on 2nd and 7th file
	game.whiteBoard[KNIGHT] = 0x42
	game.blackBoard[KNIGHT] = 0x42 << 56

	// Add knights on 3rd and 6th file
	game.whiteBoard[BISHOP] = 0x24
	game.blackBoard[BISHOP] = 0x24 << 56

	// Add pawns on 2nd and 7th file
	game.whiteBoard[PAWN] = 0xFF << 8
	game.blackBoard[PAWN] = 0xFF << 48
}
