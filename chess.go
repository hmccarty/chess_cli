package main

import (
	//"fmt"
	"errors"
)

type GetMoveSubset func(uint64, Color) uint64

const ASCII_ROW_OFFSET = 49
const ASCII_COL_OFFSET = 96

const NOT_A_FILE = 0x7f7f7f7f7f7f7f7f
const NOT_H_FILE = 0xfefefefefefefefe
const A_FILE_CORNERS = 0x8000000000000080
const H_FILE_CORNERS = 0x0100000000000001

const EIGTH_RANK = 0xFF000000000000FF

const KING_CASTLE_MASK = 0x06
const QUEEN_CASTLE_MASK = 0x30

// Enum to define board types
type Board uint8
const (
	KING Board = iota
	QUEEN
	ROOK
	BISHOP
	KNIGHT
	PAWN
	EMPTY
)

var boardToPoints = map[Board]int8 {
	KING   : 0,
	QUEEN  : 9,
	ROOK   : 5,
	BISHOP : 3,
	KNIGHT : 3,
	PAWN   : 1,
	EMPTY  : 0,
}

type Color uint8
const (
	WHITE Color = iota
	BLACK
)

var oppColor = map[Color]Color {
	WHITE : BLACK,
	BLACK : WHITE,
}

type RayDirections uint8
const (
	NORTH RayDirections = iota
	NORTH_EAST
	EAST
	SOUTH_EAST
	SOUTH
	SOUTH_WEST
	WEST
	NORTH_WEST
)

type Game struct {
	turn Color
	board [7]uint64
	color [2]uint64
	kingCastle [2]bool
	queenCastle [2]bool
	points [2]int8
	rayAttacks [64][8]uint64
}

type GameStatus uint8
const (
	IN_PLAY GameStatus = iota
	WHITE_WON
	BLACK_WON
	DRAW
)

type Flag uint8
const (
	QUIET Flag = iota
	CAPTURE
	KING_SIDE_CASTLE
	QUEEN_SIDE_CASTLE
	EP_CAPTURE
	PROMOTION
)

type Move struct {
	flag Flag
	kingCastle [2]bool
	queenCastle [2]bool
	points int8
	from uint64
	fromBoard Board
	fromColor Color
	to uint64
	toBoard Board
	toColor Color
}

type MoveList struct {
	num uint8
	move *Move
	next *MoveList
	prev *MoveList
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

	// Add knights on 2nd and 7th file
	game.color [WHITE] |= 0x42
	game.color [BLACK] |= 0x42 << 56
	game.board[KNIGHT] = 0x42 | (0x42 << 56)

	// Add bishop on 3rd and 6th file
	game.color [WHITE] |= 0x24
	game.color [BLACK] |= 0x24 << 56
	game.board[BISHOP] = 0x24 | (0x24 << 56)

	// Add pawns on 2nd and 7th file
	game.color [WHITE] |= 0xFF << 8
	game.color [BLACK] |= 0xFF << 48
	game.board[PAWN] = (0xFF << 8) | (0xFF << 48)

	// Track all open squares
	game.board[EMPTY] = game.FindEmptySpaces()

	// Calculate ray attacks
	// TODO: Find a more elegant approach to ray-move calculation
	for i, _ := range game.rayAttacks {
		row := i / 8
		col := i % 8  
		
		// Calculate north ray attacks
		for j := 8; j > row; j-- {
			game.rayAttacks[i][NORTH] = moveNorth((1 << i) | game.rayAttacks[i][NORTH])
		}

		// Calculate north-east ray attacks
		for j := 8; j > row; j-- {
			game.rayAttacks[i][NORTH_EAST] = moveNEast((1 << i) | game.rayAttacks[i][NORTH_EAST])
		}

		// Calculate east ray attacks
		for j := 0; j < col; j++ {
			game.rayAttacks[i][EAST] = moveEast((1 << i) | game.rayAttacks[i][EAST])
		}

		// Calculate south-east ray attacks
		for j := row; j > 0; j-- {
			game.rayAttacks[i][SOUTH_EAST] = moveSEast((1 << i) | game.rayAttacks[i][SOUTH_EAST])
		}

		// Calculate south ray attacks
		for j := row; j > 0; j-- {
			game.rayAttacks[i][SOUTH] = moveSouth((1 << i) | game.rayAttacks[i][SOUTH])
		}

		// Calculate south-west ray attacks
		for j := row; j > 0; j-- {
			game.rayAttacks[i][SOUTH_WEST] = moveSWest((1 << i) | game.rayAttacks[i][SOUTH_WEST])
		}

		// Calculate west ray attacks
		for j := col; j < 8; j++ {
			game.rayAttacks[i][WEST] = moveWest((1 << i) | game.rayAttacks[i][WEST])
		}

		// Calculate north-west ray attacks
		for j := 8; j > row; j-- {
			game.rayAttacks[i][NORTH_WEST] = moveNWest((1 << i) | game.rayAttacks[i][NORTH_WEST])
		}
	}

	game.kingCastle[WHITE] = true
	game.kingCastle[BLACK] = true
	game.queenCastle[WHITE] = true
	game.queenCastle[BLACK] = true
	game.turn = WHITE
}

func (game *Game) TranslateCommand(cmd string) (uint8, uint8) {
	var cmdData []byte = []byte(cmd)
	var startCol uint8 = 8 - (cmdData[0] - ASCII_COL_OFFSET)
	var startRow uint8 = cmdData[1] - ASCII_ROW_OFFSET
	var endCol uint8 = 8 - (cmdData[2] - ASCII_COL_OFFSET)
	var endRow uint8 = cmdData[3] - ASCII_ROW_OFFSET
	var fromSqr uint8 = (startRow * 8) + startCol
	var toSqr uint8 = (endRow * 8) + endCol
	return fromSqr, toSqr
}

func TranslateMove(move *Move) string {
	var fromSqr uint8 = bitScanForward(move.from)
	var toSqr uint8 = bitScanForward(move.to)
	var startRow uint8 = fromSqr / 8
	var startCol uint8 = fromSqr % 8
	var endRow uint8 = toSqr / 8
	var endCol uint8 = toSqr % 8
	return (string((8 - startCol) + ASCII_COL_OFFSET) +
	        string(startRow + ASCII_ROW_OFFSET) + 
			string((8 - endCol) + ASCII_COL_OFFSET) +
			string(endRow + ASCII_ROW_OFFSET))
}

func (game *Game) ProcessMove(fromSqr uint8, toSqr uint8) (*Move, error) {
	var from uint64 = 1 << fromSqr
	var to uint64 = 1 << toSqr

	var fromBoard Board = game.FindBoard(from)
	var fromColor Color = game.FindColor(from)
	var toBoard Board = game.FindBoard(to)
	var toColor Color = game.FindColor(to)

	if (toBoard == EMPTY) {
		toBoard = fromBoard
		toColor = fromColor
	}

	if (fromColor != game.turn) {
		return nil, errors.New("Cannot move opponent's piece.")
	}

	var flag Flag = QUIET

	switch fromBoard {
	case KING:
		if ((to & game.GetKingMoves(from, fromColor)) == 0) {
			if (to & (game.board[KING] >> 2) != 0) {
				flag = KING_SIDE_CASTLE
			} else if (to & (game.board[KING] << 2) != 0) {
				flag = QUEEN_SIDE_CASTLE	
			} else {
				return nil, errors.New("Invalid king move.")
			}
		}
	case QUEEN:
		if ((to & game.GetQueenMoves(from, fromColor)) == 0) {
			return nil, errors.New("Invalid queen move.")
		}
	case ROOK:
		if ((to & game.GetRookMoves(from, fromColor)) == 0) {
			return nil, errors.New("Invalid rook move.")
		}
	case BISHOP:
		if ((to & game.GetBishopMoves(from, fromColor)) == 0) {
			return nil, errors.New("Invalid bishop move.")
		}
	case KNIGHT:
		if ((to & game.GetKnightMoves(from, fromColor)) == 0) {
			return nil, errors.New("Invalid knight move.")
		}
	case PAWN:
		if ((to & game.GetPawnMoves(from, fromColor)) == 0) {
			return nil, errors.New("Invalid pawn move.")
		} else if ((to & EIGTH_RANK) != 0) {
			toBoard = promptPiecePromotion()
			toColor = fromColor
		}
	case EMPTY:
		return nil, errors.New("Piece doesn't exist at square.")
	}

	var move *Move = new(Move)
	if (to & game.board[EMPTY] == 0) {
		flag = CAPTURE
		move.points = boardToPoints[toBoard]
	}

	move.flag = flag
	move.from = from
	move.fromBoard = fromBoard
	move.fromColor = fromColor
	move.to = to
	move.toBoard = toBoard
	move.toColor = toColor
	move.kingCastle[move.fromColor] = false
	move.queenCastle[move.fromColor] = false

	// TODO: Replace with branchless implementation
	if (move.fromBoard == KING) {
		if (game.kingCastle[move.fromColor] == true) {
			move.kingCastle[move.fromColor] = true
		}
		if (game.queenCastle[move.fromColor] == true) {
			move.queenCastle[move.fromColor] = true
		}
	} else if (move.fromBoard == ROOK) {
		if ((move.from & A_FILE_CORNERS) != 0) {
			if (game.queenCastle[move.fromColor] == true) {
				move.queenCastle[move.fromColor] = true
			}
		} else if ((move.from & H_FILE_CORNERS) != 0) {
			if (game.kingCastle[move.fromColor] == true) {
				move.kingCastle[move.fromColor] = true
			}
		}
	}

	game.MakeMove(move)
	var inCheck bool = game.IsKingInCheck(move.fromColor)
	game.UndoMove(move)
	
	if (inCheck) {
		return nil, errors.New("King would be in check.")
	}

	return move, nil
}

func (game *Game) MakeMove(move *Move) {
	// If not capturing any pieces
	switch move.flag {
	case QUIET:
		game.QuietMove(move)
	case CAPTURE:
		game.Capture(move)
	case KING_SIDE_CASTLE:
		game.CastleKingSide(move)
	case QUEEN_SIDE_CASTLE:
		game.CastleQueenSide(move)
	}

	if (move.kingCastle[move.fromColor]) {
		game.kingCastle[move.fromColor] = !game.kingCastle[move.fromColor]
	}

	if (move.queenCastle[move.fromColor]) {
		game.queenCastle[move.fromColor] = !game.queenCastle[move.fromColor]
	}

	game.turn = oppColor[game.turn]
	game.board[EMPTY] = game.FindEmptySpaces()
}

func (game *Game) UndoMove(move *Move) {
	move.points *= -1
	game.MakeMove(move)
	move.points *= -1
}

func (game *Game) QuietMove(move *Move) {
	game.board[move.fromBoard] ^= move.from
	game.board[move.toBoard] ^= move.to
	game.color[move.fromColor] ^= move.from
	game.color[move.toColor] ^= move.to
}

func (game *Game) Capture(move *Move) {
	// Remove attacked piece
	game.board[move.toBoard] ^= move.to
	game.color[move.toColor] ^= move.to

	// Move piece on attacking board
	game.board[move.fromBoard] ^= (move.from ^ move.to)
	game.color[move.fromColor] ^= (move.from ^ move.to)

	// Update point totals
	game.points[move.fromColor] += boardToPoints[move.toBoard]
}

func (game *Game) CastleKingSide(move *Move) {
	if (move.fromColor == WHITE) {
		game.board[KING] ^= 0x0A
		game.board[ROOK] ^= 0x05
		game.color[move.fromColor] ^= 0x0F
	} else {
		game.board[KING] ^= (0x0A << 56)
		game.board[ROOK] ^= (0x05 << 56)
		game.color[move.fromColor] ^= 0x0F << 56
	}
}

func (game *Game) CastleQueenSide(move *Move) {
	if (move.fromColor == WHITE) {
		game.board[KING] ^= 0x28
		game.board[ROOK] ^= 0x90
		game.color[move.fromColor] ^= 0xB8
	} else {
		game.board[KING] ^= (0x28 << 56)
		game.board[ROOK] ^= (0x90 << 56)
		game.color[move.fromColor] ^= 0xB8 << 56
	}
}

func (game *Game) FindEmptySpaces() uint64 {
	return ^(game.color[WHITE] | game.color[BLACK])
}

func (game *Game) FindBoard(pos uint64) Board {
	for idx, piece := range game.board {
		if ((piece & pos) != 0) {
			return Board(idx)
		}
	}
	return EMPTY
}

func (game *Game) FindColor(pos uint64) Color {
	if ((game.color[WHITE] & pos) != 0) {
		return WHITE
	} else {
		return BLACK
	}
}

func (game *Game) GetAllLegalMoves() *MoveList {
	var first *MoveList = new(MoveList)
	first.num = 1

	var pieces uint64 = 0
	var currMove *MoveList = first

	pieces = game.board[KING] & game.color[game.turn]
	currMove = game.AddMovesToList(pieces, game.turn, game.GetKingMoves, currMove)

	pieces = game.board[QUEEN] & game.color[game.turn]
	currMove = game.AddMovesToList(pieces, game.turn, game.GetQueenMoves, currMove)

	pieces = game.board[ROOK] & game.color[game.turn]
	currMove = game.AddMovesToList(pieces, game.turn, game.GetRookMoves, currMove)

	pieces = game.board[BISHOP] & game.color[game.turn]
	currMove = game.AddMovesToList(pieces, game.turn, game.GetBishopMoves, currMove)

	pieces = game.board[KNIGHT] & game.color[game.turn]
	currMove = game.AddMovesToList(pieces, game.turn, game.GetKnightMoves, currMove)

	pieces = game.board[PAWN] & game.color[game.turn]
	game.AddMovesToList(pieces, game.turn, game.GetPawnMoves, currMove)

	return first
}

func (game *Game) AddMovesToList(pieces uint64, color Color,
								 getMovesSubset GetMoveSubset,
								 moveList *MoveList) *MoveList {
	for (pieces != 0) {
		var piece uint64 = 1 << bitScanForward(pieces)
		var moves uint64 = getMovesSubset(piece, color)
		for (moves != 0) {
			var toSqr uint8 = bitScanForward(moves)
			move, err := game.ProcessMove(bitScanForward(piece), toSqr)
			if (err == nil) {
				moveList.move = move
				moveList.next = new(MoveList)
				moveList.next.prev = moveList
				moveList = moveList.next
			}
			moves ^= 1 << toSqr
		}
		pieces ^= piece
	}

	return moveList
}

func (game *Game) CanCastleKingSide(color Color) bool {
	if (!game.kingCastle[color]) {
		return false
	}

	var castleMask uint64 = KING_CASTLE_MASK
	if (color == BLACK) {
		castleMask = castleMask << 56
	}

	if ((game.board[EMPTY] & castleMask) != castleMask) {
		return false
	} else if (game.IsSqrUnderAttack(bitScanForward(castleMask), color) ||
			   game.IsSqrUnderAttack(bitScanReverse(castleMask), color)) {
		return false
	} else {
		return true
	}
}

func (game *Game) CanCastleQueenSide(color Color) bool {
	if (!game.queenCastle[color]) {
		return false
	}

	var castleMask uint64 = QUEEN_CASTLE_MASK
	if (color == BLACK) {
		castleMask = castleMask << 56
	}

	if ((game.board[EMPTY] & castleMask) != castleMask) {
		return false
	} else if (game.IsSqrUnderAttack(bitScanForward(castleMask), color) ||
			  game.IsSqrUnderAttack(bitScanReverse(castleMask), color)) {
		return false
	} else {
		return true
	}
}

func (game *Game) GetKingMoves(board uint64, color Color) uint64 {
	var moves uint64 = moveNorth(board) | moveSouth(board)
	moves |= moveEast(board) | moveWest(board)
	moves |= moveNEast(board) | moveNWest(board)
	moves |= moveSEast(board) | moveSWest(board)
	return moves & (^game.color[color])
}

func (game *Game) GetKnightMoves(board uint64, color Color) uint64 {
	var moves uint64 = moveNorth(moveNEast(board) | moveNWest(board))
	moves |= moveEast(moveNEast(board) | moveSEast(board))
	moves |= moveWest(moveNWest(board) | moveSWest(board))
	moves |= moveSouth(moveSEast(board) | moveSWest(board))
	return moves & (^game.color[color])
}

func (game *Game) GetRookMoves(board uint64, color Color) uint64 {
	return game.GetTransMoves(board) & (^game.color[color])
}

func (game *Game) GetBishopMoves(board uint64, color Color) uint64 {
	return game.GetDiagMoves(board) & (^game.color[color])
}

func (game *Game) GetQueenMoves(board uint64, color Color) uint64 {
	var moves uint64 = game.GetTransMoves(board) | game.GetDiagMoves(board)
	return moves & (^game.color[color])
}

func (game *Game) GetPawnMoves(board uint64, color Color) uint64 {
	var moves uint64 = 0
	if (color == WHITE) {
		// Check for single push
		moves |= (moveNorth(board) & game.board[EMPTY])
		// Check for double push
		moves |= ((0xFF << 24) & moveNorth(moves) & game.board[EMPTY])
		// Check for north-east attack
		moves |= (moveNEast(board) & game.color[BLACK])
		// Check for north-west attack
		moves |= (moveNWest(board) & game.color[BLACK])
	} else {
		// Check for single push
		moves |= (moveSouth(board) & game.board[EMPTY])
		// Check for double push
		moves |= ((0xFF << 32) & moveSouth(moves) & game.board[EMPTY])
		// Check for south-east attack
		moves |= (moveSEast(board) & game.color[WHITE])
		// Check for south-west attack
		moves |= (moveSWest(board) & game.color[WHITE])
	}

	return moves
}

func (game *Game) IsKingInCheck(color Color) bool {
	var king uint64 = game.board[KING] & game.color[color]
	return game.IsSqrUnderAttack(bitScanForward(king), color)
}

func (game *Game) IsSqrUnderAttack(sqr uint8, color Color) bool {
	var pos uint64 = 1 << sqr
	var attacks uint64 = 0

	attacks |= game.GetKingMoves(pos, color) & game.board[KING]
	attacks |= game.GetRookMoves(pos, color) & (game.board[ROOK] | game.board[QUEEN])
	attacks |= game.GetBishopMoves(pos, color) & (game.board[BISHOP] | game.board[QUEEN])
	attacks |= game.GetKnightMoves(pos, color) & game.board[KNIGHT]
	attacks |= game.GetPawnMoves(pos, color) & game.board[PAWN]
	
	return (attacks != 0)
}

func (game *Game) GetGameStatus() GameStatus {
	var legalMoves *MoveList = game.GetAllLegalMoves()

	// If no legal moves, checkmate
	if (legalMoves.move == nil) {
		if (game.turn == WHITE) {
			return BLACK_WON
		} else {
			return WHITE_WON 
		}
	}

	return IN_PLAY
}

func (game *Game) GetTransMoves(piece uint64) uint64 {
	var moves uint64 = 0
	for piece != 0 {
		var sqr uint8 = bitScanForward(piece)
		moves |= game.GetPosRayAttacks(sqr, NORTH)
		moves |= game.GetNegRayAttacks(sqr, EAST)
		moves |= game.GetPosRayAttacks(sqr, WEST)
		moves |= game.GetNegRayAttacks(sqr, SOUTH)
		piece ^= (1 << sqr)
	}
	return moves
}

func (game *Game) GetDiagMoves(piece uint64) uint64 {
	var moves uint64 = 0
	for piece != 0 {
		var sqr uint8 = bitScanForward(piece)
		moves |= game.GetPosRayAttacks(sqr, NORTH_EAST)
		moves |= game.GetPosRayAttacks(sqr, NORTH_WEST)
		moves |= game.GetNegRayAttacks(sqr, SOUTH_EAST)
		moves |= game.GetNegRayAttacks(sqr, SOUTH_WEST)
		piece ^= (1 << sqr)
	}
	return moves
}

func (game *Game) GetPosRayAttacks(sqr uint8, dir RayDirections) uint64 {
	var attacks uint64 = game.rayAttacks[sqr][dir]
	var blockers uint64 = attacks & (^game.board[EMPTY])
	sqr = bitScanForward(blockers | (0x8000000000000000))
	return attacks ^ game.rayAttacks[sqr][dir]
}

func (game *Game) GetNegRayAttacks(sqr uint8, dir RayDirections) uint64 {
	var attacks uint64 = game.rayAttacks[sqr][dir]
	var blockers uint64 = attacks & (^game.board[EMPTY])
	sqr = bitScanReverse(blockers | 1)
	return attacks ^ game.rayAttacks[sqr][dir]
}

func moveNWest(board uint64) uint64 {return (board << 9) & (NOT_H_FILE)}
func moveNorth(board uint64) uint64 {return board << 8}
func moveNEast(board uint64) uint64 {return (board << 7) & (NOT_A_FILE)}
func moveEast(board uint64) uint64 {return (board >> 1) & (NOT_A_FILE)}
func moveSEast(board uint64) uint64 {return (board >> 9) & (NOT_A_FILE)}
func moveSouth(board uint64) uint64 {return board >> 8}
func moveSWest(board uint64) uint64 {return (board >> 7) & (NOT_H_FILE)}
func moveWest(board uint64) uint64 {return (board << 1) & (NOT_H_FILE)}

func decodePosition(row uint8, col uint8) uint64 {
	return 0x1 << ((uint64(row) << 3) | uint64(col))
}