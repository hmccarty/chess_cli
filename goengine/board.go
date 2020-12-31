package goengine

import (
	"errors"
	"fmt"
)

type Board struct {
	piece [7]uint64
	color [2]uint64
	rayAttacks [64][8]uint64
	kingCastle [2]bool
	queenCastle [2]bool	
	enPassant uint64
}

func (board *Board) setup() {
	// Add pieces to bitboards, bitshifting math is
	// redundant but better represents positions within
	// the bitmap

	// Add kings on 5th file
	board.color[WHITE] |= 0x08
	board.color[BLACK] |= 0x08 << 56
	board.piece[KING] = 0x08 | (0x08 << 56)

	// Add queens on 4th file
	board.color[WHITE] |= 0x10
	board.color[BLACK] |= 0x10 << 56
	board.piece[QUEEN] = 0x10 | (0x10 << 56)

	// Add rooks on 1st and 8th file
	board.color[WHITE] |= 0x81
	board.color[BLACK] |= 0x81 << 56
	board.piece[ROOK] = 0x81 | (0x81 << 56)

	// Add knights on 2nd and 7th file
	board.color[WHITE] |= 0x42
	board.color[BLACK] |= 0x42 << 56
	board.piece[KNIGHT] = 0x42 | (0x42 << 56)

	// Add bishop on 3rd and 6th file
	board.color[WHITE] |= 0x24
	board.color[BLACK] |= 0x24 << 56
	board.piece[BISHOP] = 0x24 | (0x24 << 56)

	// Add pawns on 2nd and 7th file
	board.color[WHITE] |= 0xFF << 8
	board.color[BLACK] |= 0xFF << 48
	board.piece[PAWN] = (0xFF << 8) | (0xFF << 48)

	// Track all open squares
	board.piece[EMPTY] = board.getEmptySet()

	// Calculate ray attacks
	// TODO: Find a more elegant approach to ray-move calculation
	for i, _ := range board.rayAttacks {
		row := i / 8
		col := i % 8  
		
		// Calculate north ray attacks
		for j := 8; j > row; j-- {
			board.rayAttacks[i][NORTH] = moveNorth((1 << i) | board.rayAttacks[i][NORTH])
		}

		// Calculate north-east ray attacks
		for j := 8; j > row; j-- {
			board.rayAttacks[i][NORTH_EAST] = moveNEast((1 << i) | board.rayAttacks[i][NORTH_EAST])
		}

		// Calculate east ray attacks
		for j := 0; j < col; j++ {
			board.rayAttacks[i][EAST] = moveEast((1 << i) | board.rayAttacks[i][EAST])
		}

		// Calculate south-east ray attacks
		for j := row; j > 0; j-- {
			board.rayAttacks[i][SOUTH_EAST] = moveSEast((1 << i) | board.rayAttacks[i][SOUTH_EAST])
		}

		// Calculate south ray attacks
		for j := row; j > 0; j-- {
			board.rayAttacks[i][SOUTH] = moveSouth((1 << i) | board.rayAttacks[i][SOUTH])
		}

		// Calculate south-west ray attacks
		for j := row; j > 0; j-- {
			board.rayAttacks[i][SOUTH_WEST] = moveSWest((1 << i) | board.rayAttacks[i][SOUTH_WEST])
		}

		// Calculate west ray attacks
		for j := col; j < 8; j++ {
			board.rayAttacks[i][WEST] = moveWest((1 << i) | board.rayAttacks[i][WEST])
		}

		// Calculate north-west ray attacks
		for j := 8; j > row; j-- {
			board.rayAttacks[i][NORTH_WEST] = moveNWest((1 << i) | board.rayAttacks[i][NORTH_WEST])
		}
	}

	board.kingCastle[WHITE] = true
	board.kingCastle[BLACK] = true
	board.queenCastle[WHITE] = true
	board.queenCastle[BLACK] = true
}

func (board *Board) processMove(move *Move) error {
	move.from.piece = board.findPiece(move.from.board)
	move.from.color = board.findColor(move.from.board)
	move.to.piece = board.findPiece(move.to.board)
	move.to.color = board.findColor(move.to.board)

	if move.from.piece == EMPTY {
		return errors.New("Piece doesn't exist at square.")
	} else if (move.to.piece == EMPTY) {
		move.flag = QUIET
		move.to.piece = move.from.piece
		move.to.color = move.from.color
	}

	var pieceSet GetSet = board.getPieceSet(move.from.piece)
	var setRange uint64 = pieceSet(move.from.board, move.from.color)
	if (move.to.board & setRange) == 0 {
		if move.from.piece == KING {
			if (move.to.board & (board.piece[KING] >> 2)) != 0 {
				move.flag = KING_SIDE_CASTLE
			} else if (move.to.board & (board.piece[KING] << 2)) != 0 {
				move.flag = QUEEN_SIDE_CASTLE
			} else {
				return errors.New("Move does not exist in piece range.")
			}
		} else if move.from.piece == PAWN {
			if (board.getPawnAttackSet(move.from.board, move.from.color) & move.from.enPassant) != 0 {
				move.flag = EP_CAPTURE
				move.to.points[move.from.color] = pieceToPoints[PAWN]
			}
		} else {
			return errors.New("Move does not exist in piece range.")
		}
	}

	if move.from.piece == PAWN {
		if (move.to.board & EIGTH_RANK) != 0 {
			move.flag = PROMOTION
		} else {
			if (move.from.color == WHITE) {
				if (moveNorth(moveNorth(move.from.board)) & move.to.board) != 0 {
					move.to.enPassant = moveSouth(move.to.board)
				}
			} else {
				if (moveSouth(moveSouth(move.from.board)) & move.to.board) != 0 {
					move.to.enPassant = moveNorth(move.to.board)
				}
			}
		}
	}	

	if (move.to.board & board.piece[EMPTY]) == 0 {
		move.flag = CAPTURE
		move.to.points[move.from.color] = pieceToPoints[move.to.piece]
	}

	// TODO: Replace with branchless implementation
	if (move.from.piece == KING) {
		move.to.kingCastle[move.from.color] = false
		move.to.queenCastle[move.from.color] = false
	} else if (move.from.piece == ROOK) {
		move.to.kingCastle[move.from.color] = false
		move.to.queenCastle[move.from.color] = false
	}

	return nil
}

func (board *Board) quietMove(move *Move) {
	board.piece[move.from.piece] ^= move.from.board
	board.piece[move.to.piece] ^= move.to.board
	board.color[move.from.color] ^= move.from.board
	board.color[move.to.color] ^= move.to.board

	board.piece[EMPTY] = board.getEmptySet()
}

func (board *Board) capture(move *Move) {
	// Remove attacked piece
	board.piece[move.to.piece] ^= move.to.board
	board.color[move.to.color] ^= move.to.board

	// Move piece on attacking board
	board.piece[move.from.piece] ^= (move.from.board ^ move.to.board)
	board.color[move.from.color] ^= (move.from.board ^ move.to.board)

	board.piece[EMPTY] = board.getEmptySet()
}

func (board *Board) castleKingSide(move *Move) {
	if (move.from.color == WHITE) {
		board.piece[KING] ^= 0x0A
		board.piece[ROOK] ^= 0x05
		board.color[move.from.color] ^= 0x0F
	} else {
		board.piece[KING] ^= (0x0A << 56)
		board.piece[ROOK] ^= (0x05 << 56)
		board.color[move.from.color] ^= 0x0F << 56
	}

	board.piece[EMPTY] = board.getEmptySet()
}

func (board *Board) castleQueenSide(move *Move) {
	if (move.from.color == WHITE) {
		board.piece[KING] ^= 0x28
		board.piece[ROOK] ^= 0x90
		board.color[move.from.color] ^= 0xB8
	} else {
		board.piece[KING] ^= (0x28 << 56)
		board.piece[ROOK] ^= (0x90 << 56)
		board.color[move.from.color] ^= 0xB8 << 56
	}

	board.piece[EMPTY] = board.getEmptySet()
}

func (board *Board) updateCastleRights(move *Move) {
	board.kingCastle[WHITE] = move.to.kingCastle[WHITE]
	board.kingCastle[BLACK] = move.to.kingCastle[BLACK]
}

func (board *Board) findPiece(pos uint64) Piece {
	for idx, piece := range board.piece {
		if ((piece & pos) != 0) {
			return Piece(idx)
		}
	}
	return EMPTY
}

func (board *Board) findColor(pos uint64) Color {
	if ((board.color[WHITE] & pos) != 0) {
		return WHITE
	} else {
		return BLACK
	}
}

func (board *Board) setFENString(fen string) {
	// TODO
}

func (board *Board) getFENBoard() string {
	var mailbox [8][8]uint8 = board.getMailbox()
	var fen string = ""

	var empty int = 0
	for i := 7; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			var piece Piece = Piece(mailbox[i][j] & 0x0F)
			if (piece == EMPTY) {
				empty += 1
				continue
			} else if (empty > 0) {
				fen += fmt.Sprintf("%d", empty)
				empty = 0
			}
		
			if ((mailbox[i][j] & WHITE_MASK) != 0) {
				fen += pieceToString[WHITE][piece]
			} else {
				fen += pieceToString[BLACK][piece]
			}
		}

		if (empty > 0) {
			fen += fmt.Sprintf("%d", empty)
			empty = 0
		}

		if (i > 0) {
			fen += "/"
		}
	}
	return fen
}

func (board *Board) getMailbox() [8][8]uint8 {
	var mailbox [8][8]uint8
	for i := 0; i < 64; i++ {
		var sqr uint64 = 1 << i
		var row int = i / 8
		var col int = 7 - (i % 8)
		for j := 0; j < 7; j++ {
			if ((board.piece[j] & sqr) != 0) {
				mailbox[row][col] = uint8(j)
				if ((board.color[WHITE] & sqr) != 0) {
					mailbox[row][col] |= WHITE_MASK
				} else {
					mailbox[row][col] |= BLACK_MASK
				}
				break
			}
		} 
	}
	return mailbox
}

func (board *Board) getPieces(piece Piece, color Color) uint64 {
	return board.piece[piece] & board.color[color]
}

func (board *Board) getEmptySet() uint64 {
	return ^(board.color[WHITE] | board.color[BLACK])
}

func (board *Board) getKingSet(piece uint64, color Color) uint64 {
	var moves uint64 = moveNorth(piece) | moveSouth(piece)
	moves |= moveEast(piece) | moveWest(piece)
	moves |= moveNEast(piece) | moveNWest(piece)
	moves |= moveSEast(piece) | moveSWest(piece)
	return moves & (^board.color[color])
}

func (board *Board) getKnightSet(piece uint64, color Color) uint64 {
	var moves uint64 = moveNorth(moveNEast(piece) | moveNWest(piece))
	moves |= moveEast(moveNEast(piece) | moveSEast(piece))
	moves |= moveWest(moveNWest(piece) | moveSWest(piece))
	moves |= moveSouth(moveSEast(piece) | moveSWest(piece))
	return moves & (^board.color[color])
}

func (board *Board) getRookSet(piece uint64, color Color) uint64 {
	return board.getTransSet(piece) & (^board.color[color])
}

func (board *Board) getBishopSet(piece uint64, color Color) uint64 {
	return board.getDiagSet(piece) & (^board.color[color])
}

func (board *Board) getQueenSet(piece uint64, color Color) uint64 {
	var moves uint64 = board.getTransSet(piece) | board.getDiagSet(piece)
	return moves & (^board.color[color])
}

func (board *Board) getPawnSet(piece uint64, color Color) uint64 {
	var moves uint64 = 0
	if (color == WHITE) {
		// Check for single push
		moves |= (moveNorth(piece) & board.piece[EMPTY])
		// Check for double push
		moves |= ((0xFF << 24) & moveNorth(moves) & board.piece[EMPTY])
	} else {
		// Check for single push
		moves |= (moveSouth(piece) & board.piece[EMPTY])
		// Check for double push
		moves |= ((0xFF << 32) & moveSouth(moves) & board.piece[EMPTY])
	}
	moves |= board.getPawnAttackSet(piece, color)
	return moves
}

func (board *Board) getPawnAttackSet(piece uint64, color Color) uint64 {
	var moves uint64 = 0
	if (color == WHITE) {
		// Check for north-east attack
		moves |= (moveNEast(piece) & board.color[BLACK])
		// Check for north-west attack
		moves |= (moveNWest(piece) & board.color[BLACK])
	} else {
		// Check for south-east attack
		moves |= (moveSEast(piece) & board.color[WHITE])
		// Check for south-west attack
		moves |= (moveSWest(piece) & board.color[WHITE])
	}
	return moves
}

func (board *Board) getPieceSet(piece Piece) GetSet {
	switch (piece) {
		case KING:
			return board.getKingSet
		case QUEEN:
			return board.getQueenSet
		case ROOK:
			return board.getRookSet
		case BISHOP:
			return board.getBishopSet
		case KNIGHT:
			return board.getKnightSet
		case PAWN:
			return board.getPawnSet
		default:
			return nil
	}
}

func (board *Board) isKingInCheck(color Color) bool {
	var king uint64 = board.piece[KING] & board.color[color]
	return board.isSqrUnderAttack(bitScanForward(king), color)
}

func (board *Board) canCastleKingSide(color Color) bool {
	if (!board.kingCastle[color]) {
		return false
	}

	var castleMask uint64 = KING_CASTLE_MASK
	if (color == BLACK) {
		castleMask = castleMask << 56
	}

	if ((board.piece[EMPTY] & castleMask) != castleMask) {
		return false
	} else if (board.isSqrUnderAttack(bitScanForward(castleMask), color) ||
			   board.isSqrUnderAttack(bitScanReverse(castleMask), color)) {
		return false
	} else {
		return true
	}
}

func (board *Board) canCastleQueenSide(color Color) bool {
	if (!board.queenCastle[color]) {
		return false
	}

	var castleMask uint64 = QUEEN_CASTLE_MASK
	if (color == BLACK) {
		castleMask = castleMask << 56
	}

	if ((board.piece[EMPTY] & castleMask) != castleMask) {
		return false
	} else if (board.isSqrUnderAttack(bitScanForward(castleMask), color) ||
			  board.isSqrUnderAttack(bitScanReverse(castleMask), color)) {
		return false
	} else {
		return true
	}
}

func (board *Board) isSqrUnderAttack(sqr uint8, color Color) bool {
	var pos uint64 = 1 << sqr
	var attacks uint64 = 0

	attacks |= board.getKingSet(pos, color) & board.piece[KING]
	attacks |= board.getRookSet(pos, color) & (board.piece[ROOK] | board.piece[QUEEN])
	attacks |= board.getBishopSet(pos, color) & (board.piece[BISHOP] | board.piece[QUEEN])
	attacks |= board.getKnightSet(pos, color) & board.piece[KNIGHT]
	attacks |= board.getPawnSet(pos, color) & board.piece[PAWN]
	
	return (attacks != 0)
}

func (board *Board) getTransSet(piece uint64) uint64 {
	var moves uint64 = 0
	for piece != 0 {
		var sqr uint8 = bitScanForward(piece)
		moves |= board.getPosRayAttacks(sqr, NORTH)
		moves |= board.getNegRayAttacks(sqr, EAST)
		moves |= board.getPosRayAttacks(sqr, WEST)
		moves |= board.getNegRayAttacks(sqr, SOUTH)
		piece ^= (1 << sqr)
	}
	return moves
}

func (board *Board) getDiagSet(piece uint64) uint64 {
	var moves uint64 = 0
	for piece != 0 {
		var sqr uint8 = bitScanForward(piece)
		moves |= board.getPosRayAttacks(sqr, NORTH_EAST)
		moves |= board.getPosRayAttacks(sqr, NORTH_WEST)
		moves |= board.getNegRayAttacks(sqr, SOUTH_EAST)
		moves |= board.getNegRayAttacks(sqr, SOUTH_WEST)
		piece ^= (1 << sqr)
	}
	return moves
}

func (board *Board) getPosRayAttacks(sqr uint8, dir RayDirections) uint64 {
	var attacks uint64 = board.rayAttacks[sqr][dir]
	var blockers uint64 = attacks & (^board.piece[EMPTY])
	sqr = bitScanForward(blockers | (0x8000000000000000))
	return attacks ^ board.rayAttacks[sqr][dir]
}

func (board *Board) getNegRayAttacks(sqr uint8, dir RayDirections) uint64 {
	var attacks uint64 = board.rayAttacks[sqr][dir]
	var blockers uint64 = attacks & (^board.piece[EMPTY])
	sqr = bitScanReverse(blockers | 1)
	return attacks ^ board.rayAttacks[sqr][dir]
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

type GetSet func(uint64, Color) uint64

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

// Enum to define board types
type Piece uint8
const (
	KING Piece = iota
	QUEEN
	ROOK
	BISHOP
	KNIGHT
	PAWN
	EMPTY
)

var pieceToPoints = map[Piece]int8 {
	KING   : 0,
	QUEEN  : 9,
	ROOK   : 5,
	BISHOP : 3,
	KNIGHT : 3,
	PAWN   : 1,
	EMPTY  : 0,
}

var pieceToString = [2][7]string{{"K", "Q", "R", "B", "N", "P", "X"},
								 {"k", "q", "r", "b", "n", "p", "X"},}

type Color uint8
const (
	WHITE Color = iota
	BLACK
)

var oppColor = map[Color]Color {
	WHITE : BLACK,
	BLACK : WHITE,
}

const ASCII_ROW_OFFSET = 49
const ASCII_COL_OFFSET = 96

const NOT_A_FILE = 0x7f7f7f7f7f7f7f7f
const NOT_H_FILE = 0xfefefefefefefefe
const A_FILE_CORNERS = 0x8000000000000080
const H_FILE_CORNERS = 0x0100000000000001

const EIGTH_RANK = 0xFF000000000000FF

const KING_CASTLE_MASK = 0x06
const QUEEN_CASTLE_MASK = 0x30

const WHITE_MASK = 0x10
const BLACK_MASK = 0x20