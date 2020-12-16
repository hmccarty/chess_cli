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
	hasMoved bool
}

type Move struct {
	data string
	startRow int
	startCol int
	endRow int
	endCol int
	nextMove *Move
}

func createBackRank(color Color, row int) [8]*Piece {
	backRank := [8]*Piece {&Piece{ROOK, color, row, 0, false},
						   &Piece{KNIGHT, color, row, 1, false},
						   &Piece{BISHOP, color, row, 2, false},
						   &Piece{QUEEN, color, row, 3, false},
						   &Piece{KING, color, row, 4, false},
						   &Piece{BISHOP, color, row, 5, false},
						   &Piece{KNIGHT, color, row, 6, false},
						   &Piece{ROOK, color, row, 7, false}}
	return backRank
}

func createPawnRank(color Color, row int) [8]*Piece {
	pawnRank := [8]*Piece{}
	for i := 0; i < 8; i++ {
		pawnRank[i] = &Piece{PAWN, color, row, i, false}
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

	board := [8][8]*Piece{createBackRank(frontColor, 0),
						  createPawnRank(frontColor, 1),
				 	 	  createEmptyRank(),
				 	 	  createEmptyRank(),
				 	 	  createEmptyRank(),
						  createEmptyRank(),
						  createPawnRank(backColor, 6),
						  createBackRank(backColor, 7),}
	return board
}

func (game *Game) Setup() {
	game.board = createBoard()
	game.whiteKing = game.board[0][4]
	game.blackKing = game.board[7][4]
	game.turnColor = WHITE
}

func (game *Game) AddNewMove(moveString string) {
	moveData := []byte(moveString)

	startCol := int(moveData[0] - ASCII_COL_OFFSET)
	startRow := int(moveData[1] - ASCII_ROW_OFFSET)

	endCol := int(moveData[2] - ASCII_COL_OFFSET)
	endRow := int(moveData[3] - ASCII_ROW_OFFSET)

	move := Move{moveString, startRow, startCol, endRow, endCol, nil}

	if game.IsMoveValid(&move) {
		if game.lastMove == nil {
			game.moves = &move
			game.lastMove = &move
		} else {
			game.lastMove.nextMove = &move
			game.lastMove = game.lastMove.nextMove
		}
		game.numMoves += 1
	}
}

func (game *Game) IsMoveValid(move *Move) bool {
	piece := game.board[move.startRow][move.startCol]
	deadPiece := game.board[move.endRow][move.endCol]

	if (piece == nil) ||
	   (deadPiece != nil && piece.color == deadPiece.color) ||
	   (!IsWithinRange(game.board, piece, move.endRow, move.endCol)) {
		return false
	}

	game.board[move.startRow][move.startCol] = nil
	game.board[move.endRow][move.endCol] = piece
	
	piece.row = move.endRow
	piece.col = move.endCol
	if !piece.hasMoved {
		piece.hasMoved = true
	}

	return true
}

func IsWithinRange(board [8][8]*Piece, piece *Piece,
				   endRow int, endCol int) bool {
	if (piece.row > 7 || piece.row < 0) ||
	   (piece.col > 7 || piece.col < 0) ||
	   (endRow > 7 || endRow < 0) ||
	   (endCol > 7 || endCol < 0) {
		return false
	}

	rowDiff := endRow - piece.row
	colDiff := endCol - piece.col

	if (rowDiff == 0 && colDiff == 0) {
		return false
	}

	trans := (colDiff == 0) || (rowDiff == 0)
	diag := abs(colDiff) == abs(rowDiff)

	switch piece.class {
	case KING:
		return (abs(rowDiff) <= 1) && (abs(colDiff) <= 1)
	case QUEEN:
		return (trans || diag) && !PathHasObstacle(board, piece.row, piece.col,
											       endRow, endCol)
	case ROOK:
		return trans && !PathHasObstacle(board, piece.row, piece.col,
										 endRow, endCol)
	case BISHOP:
		return diag && !PathHasObstacle(board, piece.row, piece.col,
									    endRow, endCol)
	case KNIGHT:
		return ((abs(rowDiff) == 1) && (abs(colDiff) == 2)) ||
			   ((abs(rowDiff) == 2) && (abs(colDiff) == 1))
	case PAWN:
		if (colDiff > 1 || colDiff < -1) {
			return false
		}

		maxDist := 1
		if !piece.hasMoved && colDiff == 0 {
			maxDist = 2
		}

		return ((piece.color == WHITE && rowDiff <= maxDist && rowDiff > 0) ||
				(piece.color == BLACK && rowDiff >= -maxDist && rowDiff < 0)) &&
				!PathHasObstacle(board, piece.row, piece.col, endRow, endCol)
	}

	return true
}

func PathHasObstacle(board [8][8]*Piece, row int, col int,
					 endRow int, endCol int) bool {
	for (row != endRow || col != endCol) {
		if (board[row][col] != nil) {
			return true
		}
		
		row += copysign(1, endRow - row)
		col += copysign(1, endCol - col)
	}
	return false
}