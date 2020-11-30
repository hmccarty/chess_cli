package main

import (
	"fmt"
)

const ASCII_ROW_OFFSET = 49
const ASCII_COL_OFFSET = 97

// Enum to define piece classes (or types)
type Class int
const (
	KING Class = iota
	QUEEN
	ROOK
	BISHOP
	KNIGHT
	PAWN
)

// Enum to define piece colors
type Color int
const (
	WHITE Color = iota
	BLACK
)

var pieceToChar = map[Class]string{
	KING : "K",
	QUEEN: "Q",
	ROOK: "R",
	BISHOP: "B",
	KNIGHT: "N",
	PAWN: "p",
}

type Game struct {
	id string
	board [8][8]*Piece
	whiteKing *Piece
	blackKing *Piece
	moves *Move
	lastMove *Move
	numMoves int
	turnColor Color
	inCheck bool
}

type Piece struct {
	class Class
	color Color
	row int
	col int
}

type Move struct {
	data string
	nextMove *Move
}

func createBackRank(color Color, row int) [8]*Piece {
	backRank := [8]*Piece {&Piece{ROOK, color, row, 0},
						   &Piece{KNIGHT, color, row, 1},
						   &Piece{BISHOP, color, row, 2},
						   &Piece{QUEEN, color, row, 3},
						   &Piece{KING, color, row, 4},
						   &Piece{BISHOP, color, row, 5},
						   &Piece{KNIGHT, color, row, 6},
						   &Piece{ROOK, color, row, 7}}
	return backRank
}

func createPawnRank(color Color, row int) [8]*Piece {
	pawnRank := [8]*Piece{}
	for i := 0; i < 8; i++ {
		pawnRank[i] = &Piece{PAWN, color, row, i}
	}
	return pawnRank
}

func createEmptyRank() [8]*Piece {
	emptyRank := [8]*Piece{}
	return emptyRank
}

// Creates a board representation with the correct piece orientations
func createBoard() [8][8]*Piece {
	// TODO: Allow for board to face either direction
	backColor := BLACK
	frontColor := WHITE

	board := [8][8]*Piece{createBackRank(backColor, 7),
				          createPawnRank(backColor, 6),
				 	 	  createEmptyRank(),
				 	 	  createEmptyRank(),
				 	 	  createEmptyRank(),
				 		  createEmptyRank(),
				 		  createPawnRank(frontColor, 1),
				 		  createBackRank(frontColor, 0),}
	return board
}

func (game *Game) Setup() {
	game.board = createBoard()
	game.whiteKing = game.board[0][4]
	game.blackKing = game.board[7][4]
	game.turnColor = WHITE
}

func (game *Game) AddNewMove(move string) {
	fmt.Println("Checking new move...")
	if game.IsMoveValid(move) {
		fmt.Println("New move valid")
		newMove := Move{data : move}
		if game.lastMove == nil {
			game.moves = &newMove
			game.lastMove = &newMove
		} else {
			game.lastMove.nextMove = &newMove
			game.lastMove = game.lastMove.nextMove
		}
		game.numMoves += 1
	}
}

func (game *Game) IsMoveValid(move string) bool {
	moveData := []byte(move)

	startCol := int(moveData[0] - ASCII_COL_OFFSET)
	startRow := int(7 - (moveData[1] - ASCII_ROW_OFFSET))

	endCol := int(moveData[2] - ASCII_COL_OFFSET)
	endRow := int(7 - (moveData[3] - ASCII_ROW_OFFSET))

	deadPiece := game.board[endRow][endCol]
	piece := game.board[startRow][startCol]

	if (deadPiece != nil && piece.color == deadPiece.color) ||
	   (!IsWithinRange(piece, startRow, startCol, endRow, endCol)) ||
	   (deadPiece.class == KING) {
		return false
	}

	game.board[startRow][startCol] = nil
	game.board[endRow][endCol] = piece

	return true
}

func IsWithinRange(piece *Piece, startRow int, startCol int,
				   endRow int, endCol int) bool {
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

	trans := (colDiff == 0) || (rowDiff == 0)
	diag := abs(colDiff) == abs(rowDiff)

	switch piece.class {
	case KING:
		return (abs(rowDiff) <= 1) && (abs(colDiff) <= 1)
	case QUEEN:
		return trans || diag
	case ROOK:
		return trans
	case BISHOP:
		return diag
	case KNIGHT:
		return ((abs(rowDiff) == 1) && (abs(colDiff) == 2)) ||
			   ((abs(rowDiff) == 2) && (abs(colDiff) == 1))
	case PAWN:
		if (colDiff > 1 || colDiff < -1) {
			return false
		} 
		return (piece.color == WHITE && rowDiff == 1) ||
		       (piece.color == BLACK && rowDiff == -1)
	}

	return true
}