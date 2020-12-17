package main

import (
	"fmt"
)

const ASCII_ROW_OFFSET = 49
const ASCII_COL_OFFSET = 96

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
	emptyBoard uint64
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

func (game *Game) ProcessMove(move string) uint64, uint64, error {
	moveData := []byte(moveString)
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
		*board = quietMove(from, to, *board)
	}
	game.emptyBoard = game.FindEmptySpaces()
}

func (game *Game) FindEmptySpaces() uint64 {
	var empty uint64 = 0
	for _, piece := range game.whiteBoard {
		empty |= piece
	}
	for _, piece := range game.blackBoard {
		empty |= piece
	}
	return ^empty
}

func (game *Game) FindBoard(pos uint64) *uint64 {
	for idx, piece := range game.whiteBoard {
		if ((piece & pos) != 0) {
			return &(game.whiteBoard[idx])
		}
	}
	for idx, piece := range game.blackBoard {
		if ((piece & pos) != 0) {
			return &(game.blackBoard[idx])
		}
	}
	return nil
}

func moveNWest(piece int, board [6]uint64) uint64 {return board[piece] << 9}
func moveNorth(piece int, board [6]uint64) uint64 {return board[piece] << 8}
func moveNEast(piece int, board [6]uint64) uint64 {return board[piece] << 7}
func moveEast(piece int, board [6]uint64) uint64 {return board[piece] >> 1}
func moveSEast(piece int, board [6]uint64) uint64 {return board[piece] >> 9}
func moveSouth(piece int, board [6]uint64) uint64 {return board[piece] >> 8}
func moveSWest(piece int, board [6]uint64) uint64 {return board[piece] >> 7}
func moveWest(piece int, board [6]uint64) uint64 {return board[piece] << 1}

func quietMove(from uint64, to uint64, board uint64) uint64 {
	return board ^ (from ^ to)
}

func decodePosition(row uint8, col uint8) uint64 {
	fmt.Println(row)
	// fmt.Println(col)
	// fmt.Println((8 * row) + col)
	return 0x1 << ((uint64(row) << 3) | uint64(col))
}