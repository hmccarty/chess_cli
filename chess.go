package main

import (
	"strings"
	"math"
)

const ASCII_ROW_OFFSET = 49
const ASCII_COL_OFFSET = 97

const WHITE = 0x00
const BLACK = 0x80

const EMPTY = 0x00
const KING = 0x01
const QUEEN = 0x02
const ROOK = 0x03
const BISHOP = 0x04
const KNIGHT = 0x05
const PAWN = 0x06

type Game struct {
	id string
	board [8][8]byte
	whiteKingPos [2]int
	blackKingPos [2]int
	moves *Move
	lastMove *Move
	numMoves int
	turnColor byte
	inCheck bool
}

type Move struct {
	data string
	nextMove *Move
}

var pieceToChar = map[byte]string{
	EMPTY: "x",
	KING : "K",
	QUEEN: "Q",
	ROOK: "R",
	BISHOP: "B",
	KNIGHT: "N",
	PAWN: "p",
}

func createBackRank(c byte) [8]byte {
	backRank := [8]byte {ROOK | c, KNIGHT | c, BISHOP | c, QUEEN | c, KING | c, BISHOP | c, KNIGHT | c, ROOK | c}
	return backRank
}

func createPawnRank(c byte) [8]byte {
	pawnRank := [8]byte {PAWN | c, PAWN | c, PAWN | c, PAWN | c, PAWN | c, PAWN | c, PAWN | c, PAWN | c}
	return pawnRank
}

func createEmptyRank() [8]byte {
	emptyRank := [8]byte {}
	return emptyRank
}

func createBoard(whiteFront bool) [8][8]byte {
	var frontColor byte
	var backColor byte

	if whiteFront {
		frontColor = WHITE
		backColor = BLACK
	} else {
		frontColor = BLACK
		backColor = WHITE
	}

	board := [8][8]byte{createBackRank(backColor),
				        createPawnRank(backColor),
				 		createEmptyRank(),
				 		createEmptyRank(),
				 		createEmptyRank(),
				 		createEmptyRank(),
				 		createPawnRank(frontColor),
				 		createBackRank(frontColor),}
	return board
}

func (game *Game) Setup(whiteFront bool) {
	(*game).board = createBoard(whiteFront)
	if whiteFront {
		(*game).whiteKingPos = [2][2]int{0,4,}
		(*game).blackKingPos = [2][2]int{7,4,}
	} else {
		(*game).whiteKingPos = [2][2]int{7,3,}
		(*game).blackKingPos = [2][2]int{0,3,}
	}
	(*game).turnColor = WHITE
}

func (game *Game) AddNewMove(move string) {
	if (*game).IsMoveValid(move) {
		newMove := Move{data : move}
		if (*game).lastMove == nil {
			(*game).moves = &newMove
			(*game).lastMove = &newMove
		} else {
			(*game).(*lastMove).nextMove = &newMove
			(*game).lastMove = (*game).(*lastMove).nextMove
		}
		game.numMoves += 1
	}
}

func (game *Game) IsMoveValid(move string) bool {
	moveData := []byte(move)

	startCol := moveData[0] - ASCII_COL_OFFSET
	startRow := 7 - (moveData[1] - ASCII_ROW_OFFSET)

	endCol := moveData[2] - ASCII_COL_OFFSET
	endRow := 7 - (moveData[3] - ASCII_ROW_OFFSET)

	deadPiece := (*game).(*board)[endRow][endCol]
	piece := (*game).(*board)[startRow][startCol]
	(*board)[startRow][startCol] = EMPTY
	(*board)[endRow][endCol] = piece

}

func IsWithinRange(piece, startRow, startCol, endRow, endCol) bool {
	if (startRow > 7 || startRow < 0) ||
	   (startCol > 7 || startCol < 0) ||
	   (endRow > 7 || endRow < 0) ||
	   (endCol > 7 || endCol < 0) {
		return false
	}

	rowDiff := endRow - startRow
	colDiff := endCol - startCol

	if (rowDiff == 0 && colDiff == 0) {
		return false
	}

	isPieceWhite = (piece & WHITE) == WHITE 
	trans := (math.Abs(colDiff) == 0) || (math.Abs(rowDiff == 0))
	diag := (math.Abs(colDiff) == math.Abs(rowDiff)

	switch piece {
	case EMPTY:
		return false
	case KING:
		return (math.Abs(rowDiff) <= 1) && (math.Abs(colDiff) <= 1)
	case QUEEN:
		return trans || diag
	case ROOK:
		return trans
	case BISHOP:
		return diag
	case KNIGHT:
		return ((math.Abs(rowDiff) == 1) && (math.Abs(colDiff) == 2)) ||
			   ((math.Abs(rowDiff) == 2) && (math.Abs(colDiff) == 1))
	case PAWN:
		if (colDiff > 1 || colDiff < -1) return false 
		return (isPieceWhite && rowDiff == 1) || (!isPieceWhite && rowDiff == -1)
	}
}