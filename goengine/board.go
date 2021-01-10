package goengine

import (
	"errors"
	"fmt"
)

const NOT_A_FILE = 0x7f7f7f7f7f7f7f7f
const NOT_H_FILE = 0xfefefefefefefefe
const A_FILE_CORNERS = 0x8000000000000080
const H_FILE_CORNERS = 0x0100000000000001

const EIGTH_RANK = 0xFF000000000000FF

const WHITE_MASK uint8 = 0x10
const BLACK_MASK uint8 = 0x20

const K_CASTLE_MASK uint8 = 0x06
const Q_CASTLE_MASK uint8 = 0x30

type Board struct {
	piece [7]uint64
	color [2]uint64
	castle [2]uint8
	ep uint64
}

type GetSet func(uint64, Color) uint64

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
	board.piece[EMPTY] = board.findEmptySpaces()

	// Enable castling on king and queen's side
	board.castle[WHITE] = (K_CASTLE_MASK | Q_CASTLE_MASK)
	board.castle[BLACK] = (K_CASTLE_MASK | Q_CASTLE_MASK)

	// Initialize array to track piece rays
	initRayAttacks()
}

func (board *Board) processMove(move *Move) error {
	move.target = board.findPiece(move.to)

	if (move.piece == EMPTY) {
		return errors.New("Piece does not exist at square.")
	} else if (move.target == EMPTY) {
		move.target = move.piece
	}

	switch move.piece {
	case KING:
		if ((move.flag != K_CASTLE) && (move.flag != Q_CASTLE)) &&
		   ((move.to & board.getKingSet(move.from, move.color)) == 0) {
			if (move.to & (board.piece[KING] >> 2) != 0) {
				move.flag = K_CASTLE
			} else if (move.to & (board.piece[KING] << 2) != 0) {
				move.flag = Q_CASTLE	
			} else {
				return errors.New("Invalid king move.")
			}
		}
	case QUEEN:
		if ((move.to & board.getQueenSet(move.from, move.color)) == 0) {
			return errors.New("Invalid queen move.")
		}
	case ROOK:
		if ((move.to & board.getRookSet(move.from, move.color)) == 0) {
			return errors.New("Invalid rook move.")
		}
	case BISHOP:
		if ((move.to & board.getBishopSet(move.from, move.color)) == 0) {
			return errors.New("Invalid bishop move.")
		}
	case KNIGHT:
		if ((move.to & board.getKnightSet(move.from, move.color)) == 0) {
			return errors.New("Invalid knight move.")
		}
	case PAWN:
		if ((move.to & board.getPawnSet(move.from, move.color)) == 0) {
			return errors.New("Invalid pawn move.")
		} else if (move.to & EIGTH_RANK) != 0 {
			move.flag = PROMOTION
		} else if ((move.to & board.ep) != 0) &&
				  ((move.to & board.piece[PAWN]) == 0) {
			move.flag = EP_CAPTURE
			move.points = pieceToPoints[PAWN]
		} else {
			if move.color == WHITE {
				if (moveNorth(moveNorth(move.from)) & move.to) != 0 {
					move.ep = move.to | moveSouth(move.to)
				}
			} else {
				if (moveSouth(moveSouth(move.from)) & move.to) != 0 {
					move.ep = move.to | moveNorth(move.to)
				}
			}
		}
		move.halfmove = 0
	case EMPTY:
		return errors.New("Piece doesn't exist at square.")
	}

	if move.flag == UNKNOWN {
		if (move.to & board.piece[EMPTY]) == 0 {
			move.flag = CAPTURE
			move.points = pieceToPoints[move.target]
			move.halfmove = 0
		} else {
			move.flag = QUIET
		}
	} else if move.flag == K_CASTLE {
		if !board.canCastleKingSide(move.color) {
			return errors.New("Cannot castle king side.")
		}
	} else if move.flag == Q_CASTLE {
		if !board.canCastleQueenSide(move.color) {
			return errors.New("Cannot castle queen side.")
		}
	}

	move.castle[WHITE] = board.castle[WHITE]
	move.castle[BLACK] = board.castle[BLACK]
	if (move.piece == KING) {
		move.castle[move.color] = 0
	} else if (move.piece == ROOK) {
		if ((move.from & A_FILE_CORNERS) != 0) {
			move.castle[move.color] &= (^Q_CASTLE_MASK)
		} else if ((move.from & H_FILE_CORNERS) != 0) {
			move.castle[move.color] &= (^K_CASTLE_MASK)
		}
	}

	return nil
}

func (board *Board) quietMove(move *Move) {
	board.piece[move.piece] ^= move.from
	board.piece[move.target] ^= move.to
	board.color[move.color] ^= (move.from ^ move.to)

	board.piece[EMPTY] = board.findEmptySpaces()
}

func (board *Board) capture(move *Move) {
	// Remove attacked piece
	board.piece[move.target] ^= move.to
	board.color[oppColor[move.color]] ^= move.to

	// Move piece on attacking board
	board.piece[move.piece] ^= (move.from ^ move.to)
	board.color[move.color] ^= (move.from ^ move.to)

	board.piece[EMPTY] = board.findEmptySpaces()
}

func (board *Board) epCapture(move *Move) {
	board.piece[PAWN] ^= board.ep | move.from
	board.color[move.color] ^= (move.from ^ move.to)
	board.color[oppColor[move.color]] ^= (board.ep ^ move.to)
	board.piece[EMPTY] = board.findEmptySpaces()
}

func (board *Board) castleKingSide(move *Move) {
	if (move.color == WHITE) {
		board.piece[KING] ^= 0x0A
		board.piece[ROOK] ^= 0x05
		board.color[move.color] ^= 0x0F
	} else {
		board.piece[KING] ^= (0x0A << 56)
		board.piece[ROOK] ^= (0x05 << 56)
		board.color[move.color] ^= 0x0F << 56
	}

	board.piece[EMPTY] = board.findEmptySpaces()
}

func (board *Board) castleQueenSide(move *Move) {
	if (move.color == WHITE) {
		board.piece[KING] ^= 0x28
		board.piece[ROOK] ^= 0x90
		board.color[move.color] ^= 0xB8
	} else {
		board.piece[KING] ^= (0x28 << 56)
		board.piece[ROOK] ^= (0x90 << 56)
		board.color[move.color] ^= 0xB8 << 56
	}

	board.piece[EMPTY] = board.findEmptySpaces()
}

func (board *Board) findEmptySpaces() uint64 {
	return ^(board.color[WHITE] | board.color[BLACK])
}

func (board *Board) findPiece(bb uint64) Piece {
	for i, piece := range board.piece {
		if ((piece & bb) != 0) {
			return Piece(i)
		}
	}
	return EMPTY
}

func (board *Board) findColor(bb uint64) Color {
	if ((board.color[WHITE] & bb) != 0) {
		return WHITE
	} else {
		return BLACK
	}
}

func (board *Board) clear() {
	board.piece[KING] = 0
	board.piece[QUEEN] = 0
	board.piece[ROOK] = 0
	board.piece[BISHOP] = 0
	board.piece[KNIGHT] = 0
	board.piece[PAWN] = 0
	board.color[WHITE] = 0
	board.color[BLACK] = 0
}

func (board *Board) setFENBoard(fen string) {
	board.clear()

	i := 63
	for len(fen) > 0 {
		piece := fen[0]
		fen = fen[1:]
		switch piece {
		case 'K':
			board.piece[KING] |= 1 << i
			board.color[WHITE] |= 1 << i
			i -= 1
		case 'k':
			board.piece[KING] |= 1 << i
			board.color[BLACK] |= 1 << i
			i -= 1
		case 'Q':
			board.piece[QUEEN] |= 1 << i
			board.color[WHITE] |= 1 << i
			i -= 1
		case 'q':
			board.piece[QUEEN] |= 1 << i
			board.color[BLACK] |= 1 << i
			i -= 1
		case 'R':
			board.piece[ROOK] |= 1 << i
			board.color[WHITE] |= 1 << i
			i -= 1
		case 'r':
			board.piece[ROOK] |= 1 << i
			board.color[BLACK] |= 1 << i
			i -= 1
		case 'B':
			board.piece[BISHOP] |= 1 << i
			board.color[WHITE] |= 1 << i
			i -= 1
		case 'b':
			board.piece[BISHOP] |= 1 << i
			board.color[BLACK] |= 1 << i
			i -= 1
		case 'N':
			board.piece[KNIGHT] |= 1 << i
			board.color[WHITE] |= 1 << i
			i -= 1
		case 'n':
			board.piece[KNIGHT] |= 1 << i
			board.color[BLACK] |= 1 << i
			i -= 1
		case 'P':
			board.piece[PAWN] |= 1 << i
			board.color[WHITE] |= 1 << i
			i -= 1
		case 'p':
			board.piece[PAWN] |= 1 << i
			board.color[BLACK] |= 1 << i
			i -= 1
		default:
			if (piece >= byte('1')) && (piece <= byte('8')) {
				i -= (int(piece) - (ASCII_ROW_OFFSET - 1))
			}
		}
	}

	board.piece[EMPTY] = board.findEmptySpaces()
}

func (board *Board) getFENBoard() string {
	var mailbox [8][8]uint8 = board.getMailbox()
	var fen string = ""

	var empty int = 0
	for i := 7; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			var piece Piece = Piece(mailbox[i][j] & 0x0F)
			if piece == EMPTY {
				empty += 1
				continue
			} else if (empty > 0) {
				fen += fmt.Sprintf("%d", empty)
				empty = 0
			}
		
			if (mailbox[i][j] & WHITE_MASK) != 0 {
				fen += pieceToString[WHITE][piece]
			} else {
				fen += pieceToString[BLACK][piece]
			}
		}

		if empty > 0 {
			fen += fmt.Sprintf("%d", empty)
			empty = 0
		}

		if i > 0 {
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

func (board *Board) isKingInCheck(color Color) bool {
	var king uint64 = board.piece[KING] & board.color[color]
	return board.isSqrUnderAttack(bitScanForward(king), color)
}

func (board *Board) canCastleKingSide(color Color) bool {
	if (board.castle[color] & K_CASTLE_MASK) != K_CASTLE_MASK {
		return false
	}

	var castle uint64 = uint64(K_CASTLE_MASK) << (56 * color)
	if (board.piece[EMPTY] & castle) != castle {
		return false
	} else if (board.isSqrUnderAttack(bitScanForward(castle), color) ||
			   board.isSqrUnderAttack(bitScanReverse(castle), color)) {
		return false
	} else {
		return true
	}
}

func (board *Board) canCastleQueenSide(color Color) bool {
	if (board.castle[color] & Q_CASTLE_MASK) != Q_CASTLE_MASK {
		return false
	}

	var castle uint64 = uint64(Q_CASTLE_MASK) << (56 * color)
	if ((board.piece[EMPTY] & castle) != castle) {
		return false
	} else if (board.isSqrUnderAttack(bitScanForward(castle), color) ||
			   board.isSqrUnderAttack(bitScanReverse(castle), color)) {
		return false
	} else {
		return true
	}
}

func (board *Board) isSqrUnderAttack(sqr uint8, color Color) bool {
	var piece uint64 = 1 << sqr
	var attacks uint64 = 0

	attacks |= board.getKingSet(piece, color) & board.piece[KING]
	attacks |= board.getRookSet(piece, color) &
			   (board.piece[ROOK] | board.piece[QUEEN])
	attacks |= board.getBishopSet(piece, color) &
			   (board.piece[BISHOP] | board.piece[QUEEN])
	attacks |= board.getKnightSet(piece, color) & board.piece[KNIGHT]
	attacks |= board.getPawnSet(piece, color) & board.piece[PAWN]
	
	return (attacks != 0)
}

func (board *Board) getBB(piece Piece, color Color) uint64 {
	return board.piece[piece] & board.color[color]
}

func (board *Board) getPieceSet(piece Piece, bb uint64, color Color) uint64 {
	switch piece {
	case KING:
		return board.getKingSet(bb, color)
	case QUEEN:
		return board.getQueenSet(bb, color)
	case ROOK:
		return board.getRookSet(bb, color)
	case BISHOP:
		return board.getBishopSet(bb, color)
	case KNIGHT:
		return board.getKnightSet(bb, color)
	case PAWN:
		return board.getPawnSet(bb, color)
	default:
		return 0
	}
}

func (board *Board) getKingSet(bb uint64, color Color) uint64 {
	var moves uint64 = moveNorth(bb) | moveSouth(bb)
	moves |= moveEast(bb) | moveWest(bb)
	moves |= moveNEast(bb) | moveNWest(bb)
	moves |= moveSEast(bb) | moveSWest(bb)
	return moves & (^board.color[color])
}

func (board *Board) getKnightSet(bb uint64, color Color) uint64 {
	var moves uint64 = moveNorth(moveNEast(bb) | moveNWest(bb))
	moves |= moveEast(moveNEast(bb) | moveSEast(bb))
	moves |= moveWest(moveNWest(bb) | moveSWest(bb))
	moves |= moveSouth(moveSEast(bb) | moveSWest(bb))
	return moves & (^board.color[color])
}

func (board *Board) getRookSet(bb uint64, color Color) uint64 {
	return getTransSet(bb, board.piece[EMPTY]) & (^board.color[color])
}

func (board *Board) getBishopSet(bb uint64, color Color) uint64 {
	return getDiagSet(bb, board.piece[EMPTY]) & (^board.color[color])
}

func (board *Board) getQueenSet(bb uint64, color Color) uint64 {
	var moves uint64 = getTransSet(bb, board.piece[EMPTY]) |
					   getDiagSet(bb, board.piece[EMPTY])
	return moves & (^board.color[color])
}

func (board *Board) getPawnSet(bb uint64, color Color) uint64 {
	var moves uint64 = 0
	if (color == WHITE) {
		// Check for single push
		moves |= (moveNorth(bb) & board.piece[EMPTY])
		// Check for double push
		moves |= ((0xFF << 24) & moveNorth(moves) & board.piece[EMPTY])
	} else {
		// Check for single push
		moves |= (moveSouth(bb) & board.piece[EMPTY])
		// Check for double push
		moves |= ((0xFF << 32) & moveSouth(moves) & board.piece[EMPTY])
	}
	moves |= board.getPawnAttackSet(bb, color)
	return moves
}

func (board *Board) getPawnAttackSet(bb uint64, color Color) uint64 {
	var moves uint64 = 0
	if (color == WHITE) {
		// Check for north-east attack
		moves |= moveNEast(bb) & board.color[BLACK]
		moves |= moveNEast(bb) & moveNorth(board.color[BLACK]) & board.ep
		// Check for north-west attack
		moves |= moveNWest(bb) & board.color[BLACK]
		moves |= moveNWest(bb) & moveNorth(board.color[BLACK]) & board.ep
	} else {
		// Check for south-east attack
		moves |= moveSEast(bb) & board.color[WHITE]
		moves |= moveSEast(bb) & moveSouth(board.color[WHITE]) & board.ep
		// Check for south-west attack
		moves |= moveSWest(bb) & board.color[WHITE]
		moves |= moveSWest(bb) & moveSouth(board.color[WHITE]) & board.ep
	}
	return moves
}