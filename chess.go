package main

import (
	//"fmt"
	"errors"
)

const ASCII_ROW_OFFSET = 49
const ASCII_COL_OFFSET = 96

const notAFile = 0x7f7f7f7f7f7f7f7f
const notHFile = 0xfefefefefefefefe

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

var boardToPoints = map[Board]uint8 {
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
	board [7]uint64
	color [2]uint64
	points [2]uint8
	rayAttacks [64][8]uint64
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
}

func (game *Game) ProcessMove(move string) (uint64, uint64, error) {
	var moveData []byte = []byte(move)
	var startCol uint8 = 8 - (moveData[0] - ASCII_COL_OFFSET)
	var startRow uint8 = moveData[1] - ASCII_ROW_OFFSET
	var endCol uint8 = 8 - (moveData[2] - ASCII_COL_OFFSET)
	var endRow uint8 = moveData[3] - ASCII_ROW_OFFSET
	var sqr uint8 = (startRow * 8) + startCol

	var from uint64 = decodePosition(startRow, startCol)
	var to uint64 = decodePosition(endRow, endCol)

	var board Board = game.FindBoard(from)
	var color Color = game.FindColor(from)
	switch board {
	case KING:
		if ((to & game.GetKingMoves(color)) == 0) {
			return 0, 0, errors.New("Invalid king move.")
		}
	case QUEEN:
		if ((to & game.GetQueenMoves(color)) == 0) {
			return 0, 0, errors.New("Invalid queen move.")
		}
	case ROOK:
		if ((to & game.GetRookMoves(sqr, color)) == 0) {
			return 0, 0, errors.New("Invalid rook move.")
		}
	case BISHOP:
		if ((to & game.GetBishopMoves(sqr, color)) == 0) {
			return 0, 0, errors.New("Invalid bishop move.")
		}
	case KNIGHT:
		if ((to & game.GetKnightMoves(sqr, color)) == 0) {
			return 0, 0, errors.New("Invalid knight move.")
		}
	case PAWN:
		if ((to & game.GetPawnMoves(sqr, color)) == 0) {
			return 0, 0, errors.New("Invalid pawn move.")
		}
	case EMPTY:
		return 0, 0, errors.New("Piece doesn't exist at square.")
	}
	return from, to, nil
}

func (game *Game) MakeMove(from uint64, to uint64) {
	// If not capturing any pieces
	if ((to & game.FindEmptySpaces()) != 0) {
		game.QuietMove(from, to)
	} else {
		game.Capture(from, to)
	}
	game.board[EMPTY] = game.FindEmptySpaces()
}

func (game *Game) QuietMove(from uint64, to uint64) {
	var board Board = game.FindBoard(from)
	var color Color = game.FindColor(from)
	game.board[board] ^= (from ^ to)
	game.color[color] ^= (from ^ to)
}

func (game *Game) Capture(from uint64, to uint64) {
	// Remove attacked piece
	var toBoard Board = game.FindBoard(to)
	var toColor Color = game.FindColor(to)
	game.board[toBoard] ^= to
	game.color[toColor] ^= to

	// Move piece on attacking board
	var fromBoard Board = game.FindBoard(from)
	var fromColor Color = game.FindColor(from)
	game.board[fromBoard] ^= (from ^ to)
	game.color[fromColor] ^= (from ^ to)

	// Update point totals
	game.points[fromColor] += boardToPoints[toBoard]
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

func (game *Game) GetKingMoves(color Color) uint64 {
	var king uint64 = game.board[KING] & game.color[color]
	var moves uint64 = moveNorth(king) | moveSouth(king)
	moves |= moveEast(king) | moveWest(king)
	moves |= moveNEast(king) | moveNWest(king)
	moves |= moveSEast(king) | moveSWest(king)
	return moves & (^game.color[color])
}

func (game *Game) GetKnightMoves(sqr uint8, color Color) uint64 {
	var knight uint64 = game.board[KNIGHT] & (1 << sqr)
	var moves uint64 = moveNorth(moveNEast(knight) | moveNWest(knight))
	moves |= moveEast(moveNEast(knight) | moveSEast(knight))
	moves |= moveWest(moveNWest(knight) | moveSWest(knight))
	moves |= moveSouth(moveSEast(knight) | moveSWest(knight))
	return moves & (^game.color[color])
}

func (game *Game) GetRookMoves(sqr uint8, color Color) uint64 {
	return game.GetTransMoves(sqr) & (^game.color[color])
}

func (game *Game) GetBishopMoves(sqr uint8, color Color) uint64 {
	return game.GetDiagMoves(sqr) & (^game.color[color])
}

func (game *Game) GetQueenMoves(color Color) uint64 {
	var sqr uint8 = bitScanForward(game.board[QUEEN] & game.color[color])
	var moves uint64 = game.GetTransMoves(sqr) | game.GetDiagMoves(sqr)
	return moves & (^game.color[color])
}

func (game *Game) GetPawnMoves(sqr uint8, color Color) uint64 {
	var pawn uint64 = game.board[PAWN] & (1 << sqr)
	var moves uint64 = 0
	if (color == WHITE) {
		// Check for single push
		moves |= (moveNorth(pawn) & game.board[EMPTY])
		// Check for double push
		moves |= ((0xFF << 24) & moveNorth(moves) & game.board[EMPTY])
		// Check for north-east attack
		moves |= (moveNEast(pawn) & game.color[BLACK])
		// Check for north-west attack
		moves |= (moveNWest(pawn) & game.color[BLACK])
	} else {
		// Check for single push
		moves |= (moveSouth(pawn) & game.board[EMPTY])
		// Check for double push
		moves |= ((0xFF << 32) & moveSouth(moves) & game.board[EMPTY])
		// Check for south-east attack
		moves |= (moveSEast(pawn) & game.color[WHITE])
		// Check for south-west attack
		moves |= (moveSWest(pawn) & game.color[WHITE])
	}

	return moves
}

func (game *Game) GetTransMoves(sqr uint8) uint64 {
	var moves uint64 = game.GetPosRayAttacks(sqr, NORTH)
	moves |= game.GetNegRayAttacks(sqr, EAST)
	moves |= game.GetPosRayAttacks(sqr, WEST)
	moves |= game.GetNegRayAttacks(sqr, SOUTH)
	return moves
}

func (game *Game) GetDiagMoves(sqr uint8) uint64 {
	var moves uint64 = game.GetPosRayAttacks(sqr, NORTH_EAST)
	moves |= game.GetPosRayAttacks(sqr, NORTH_WEST)
	moves |= game.GetNegRayAttacks(sqr, SOUTH_EAST)
	moves |= game.GetNegRayAttacks(sqr, SOUTH_WEST)
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

func moveNWest(board uint64) uint64 {return (board << 9) & notHFile}
func moveNorth(board uint64) uint64 {return board << 8}
func moveNEast(board uint64) uint64 {return (board << 7) & notAFile}
func moveEast(board uint64) uint64 {return (board >> 1) & notAFile}
func moveSEast(board uint64) uint64 {return (board >> 9) & notAFile}
func moveSouth(board uint64) uint64 {return board >> 8}
func moveSWest(board uint64) uint64 {return (board >> 7) & notHFile}
func moveWest(board uint64) uint64 {return (board << 1) & notHFile}

func decodePosition(row uint8, col uint8) uint64 {
	return 0x1 << ((uint64(row) << 3) | uint64(col))
}