package main

import (
	"strings"
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
	ID string
	userWhite bool
	board [8][8]byte
	moves *Move
	numMoves int
	usersTurn bool
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

func updateMoveList(game *Game, moves string) {
	moveArr := strings.Split(moves, " ")
	moveData := moveArr[len(moveArr) - 1]
	newMove := Move{data : moveData}
	appendMove(&(game.moves), &newMove)
	game.numMoves += 1
	completeMove(&game.board, moveData)
}

func appendMove(lastMove **Move, newMove *Move) {
	if *lastMove == nil {
		*lastMove = newMove
	} else if (*lastMove).nextMove == nil {
		appendMove(&(*lastMove).nextMove, newMove)
	}

	(*lastMove).nextMove = newMove
}

func completeMove(board *[8][8]byte, move string) {
	moveData := []byte(move)

	startCol := moveData[0] - ASCII_COL_OFFSET
	startRow := 7 - (moveData[1] - ASCII_ROW_OFFSET)

	endCol := moveData[2] - ASCII_COL_OFFSET
	endRow := 7 - (moveData[3] - ASCII_ROW_OFFSET)

	piece := (*board)[startRow][startCol]
	(*board)[startRow][startCol] = EMPTY
	(*board)[endRow][endCol] = piece
}