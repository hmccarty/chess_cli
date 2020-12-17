package main

import (
	//"fmt"
)

const ASCII_ROW_OFFSET = 49
const ASCII_COL_OFFSET = 96

// Enum to define piece classes (or types)
type BoardType uint8
const (
	KING BoardType = iota
	QUEEN
	ROOK
	BISHOP
	KNIGHT
	PAWN
	EMPTY
)

type Color uint8
const (
	WHITE Color = iota
	BLACK
)

type Game struct {
	board [7]uint64
	color [2]uint64
}

func (game *Game) Setup() {
	// Add pieces to bitboards, bitshifting math is
	// redundant but better represents positions within
	// the bitmap

	// Add kings on 5th file
	game.color[WHITE] |= 0x08
	game.color[BLACK] |= 0x08 << 56
	game.board[KING] = 0x08 | (0x08 << 56)

	// Add queens on 4th file
	game.color[WHITE] |= 0x10
	game.color[BLACK] |= 0x10 << 56
	game.board[QUEEN] = 0x10 | (0x10 << 56)

	// Add rooks on 1st and 8th file
	game.color [WHITE] |= 0x81
	game.color [BLACK] |= 0x81 << 56
	game.board[ROOK] = 0x81 | (0x81 << 56)

	// Add bishops on 2nd and 7th file
	game.color [WHITE] |= 0x42
	game.color [BLACK] |= 0x42 << 56
	game.board[BISHOP] = 0x42 | (0x42 << 56)

	// Add knights on 3rd and 6th file
	game.color [WHITE] |= 0x24
	game.color [BLACK] |= 0x24 << 56
	game.board[KNIGHT] = 0x24 | (0x24 << 56)

	// Add pawns on 2nd and 7th file
	game.color [WHITE] |= 0xFF << 8
	game.color [BLACK] |= 0xFF << 48
	game.board[PAWN] = (0xFF << 8) | (0xFF << 48)

	// Track all open squares
	game.board[EMPTY] = game.FindEmptySpaces()
}

func (game *Game) ProcessMove(move string) (uint64, uint64, error) {
	moveData := []byte(move)
	startCol := uint8(8 - (moveData[0] - ASCII_COL_OFFSET))
	startRow := uint8(moveData[1] - ASCII_ROW_OFFSET)
	endCol := uint8(8 - (moveData[2] - ASCII_COL_OFFSET))
	endRow := uint8(moveData[3] - ASCII_ROW_OFFSET)

	var from uint64 = decodePosition(startRow, startCol)
	var to uint64 = decodePosition(endRow, endCol)
	return from, to, nil
}

func (game *Game) MakeMove(from uint64, to uint64) {
	// If not capturing any pieces
	if ((to & game.FindEmptySpaces()) != 0) {
		var board *uint64 = game.FindBoard(from)
		var color *uint64 = &game.color[game.FindColor(from)]
		*board = quietMove(from, to, *board)
		*color = quietMove(from, to, *color)
	}
	game.board[EMPTY] = game.FindEmptySpaces()
}

func (game *Game) FindEmptySpaces() uint64 {
	return ^(game.color[WHITE] & game.color[BLACK])
}

func (game *Game) FindBoard(pos uint64) *uint64 {
	for idx, piece := range game.board {
		if ((piece & pos) != 0) {
			return &(game.board[idx])
		}
	}
	return nil
}

func (game *Game) FindColor(pos uint64) Color {
	if ((game.color[WHITE] & pos) != 0) {
		return WHITE
	} else {
		return BLACK
	}
}

func moveNWest(board uint64) uint64 {return board << 9}
func moveNorth(board uint64) uint64 {return board << 8}
func moveNEast(board uint64) uint64 {return board << 7}
func moveEast(board uint64) uint64 {return board >> 1}
func moveSEast(board uint64) uint64 {return board >> 9}
func moveSouth(board uint64) uint64 {return board >> 8}
func moveSWest(board uint64) uint64 {return board >> 7}
func moveWest(board uint64) uint64 {return board << 1}

func quietMove(from uint64, to uint64, board uint64) uint64 {
	return board ^ (from ^ to)
}

func decodePosition(row uint8, col uint8) uint64 {
	return 0x1 << ((uint64(row) << 3) | uint64(col))
}