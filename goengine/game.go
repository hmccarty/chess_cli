package goengine

import (
	"fmt"
	"errors"
	"github.com/fatih/color"
)

type Game struct {
	board *Board
	turn Color
	points [2]int8
	status GameStatus
}

type GameStatus uint8
const (
	IN_PLAY GameStatus = iota
	WHITE_WON
	BLACK_WON
	DRAW
)

func (game *Game) Setup() {
	game.board = new(Board)
	game.board.setup()
	game.turn = WHITE
}

func (game *Game) GetFENString() string {
	return game.board.getFENBoard()
}

func (game *Game) ProcessCommand(cmd string) (uint64, uint64) {
	var from uint64 = 0
	var to uint64 = 0

	if cmd == "O-O" {
		from = game.board.piece[KING] & game.board.color[game.turn]
		return from, (from >> 2)
	} else if cmd == "O-O-O" {
		from = game.board.piece[KING] & game.board.color[game.turn]
		return from, (from << 2)
	}

	var cmdData []byte = []byte(cmd)

	// If move causes check, remove 
	if cmdData[len(cmdData) - 1:][0] == byte('+') {
		cmdData = cmdData[:len(cmdData) - 1]
	}

	// Find square the piece is moving to
	var toCol uint8 = 8 - (cmdData[len(cmdData) - 2:len(cmdData) - 1][0] - ASCII_COL_OFFSET)
	var toRow uint8 = byte(cmdData[len(cmdData) - 1:][0]) - ASCII_ROW_OFFSET
	to = 1 << ((toRow * 8) + toCol)

	// Store additional info about piece
	var additionalInfo []byte

	// Determine type of piece being moved
	var fromPiece Piece
	switch cmdData[0] {
	case 'K':
		fromPiece = KING
		additionalInfo = cmdData[1:len(cmdData) - 2]
	case 'Q':
		fromPiece = QUEEN
		additionalInfo = cmdData[1:len(cmdData) - 2]
	case 'R':
		fromPiece = ROOK
		additionalInfo = cmdData[1:len(cmdData) - 2]
	case 'B':
		fromPiece = BISHOP
		additionalInfo = cmdData[1:len(cmdData) - 2]
	case 'N':
		fromPiece = KNIGHT
		additionalInfo = cmdData[1:len(cmdData) - 2]
	default:
		fromPiece = PAWN
		additionalInfo = cmdData[0:len(cmdData) - 2]
	}

	// Search all possible pieces
	// TODO: Check that piece move wouldn't cause check, otherwise look at using other pieces
	var pieces uint64 = game.board.piece[fromPiece] & game.board.color[game.turn]
	for pieces > 0 {
		var idx uint8 = bitScanForward(pieces)
		var piece uint64 = 1 << idx
		// If square piece is moving to is in piece range
		if (game.board.getPieceSet(fromPiece, piece, game.turn) & to) != 0 {
			if len(additionalInfo) == 0 {
				return piece, to
			}

			var valid bool = true
			for i := 0; i < len(additionalInfo); i++ {
				if additionalInfo[i] >= byte('a') && additionalInfo[i] <= byte('h') {
					if (8 - (idx % 8)) != (additionalInfo[i] - ASCII_COL_OFFSET) {
						valid = false
						break
					}
				} else if additionalInfo[i] >= byte('1') && additionalInfo[i]  <= byte('8') {
					if (idx / 8) != (additionalInfo[i] - ASCII_ROW_OFFSET) {
						valid = false
						break
					}
				}
			}

			if valid {
				return piece, to
			}
		}

		pieces ^= piece
	}

	return 0, 0
}

func (game *Game) CreateMove(from uint64, to uint64) (*Move, error) {
	var move *Move = new(Move)
	move.from = from
	move.to = to
	return move, game.board.processMove(move)
}

func (game *Game) CheckMove(move *Move) error {
	if (move.fromColor != game.turn) {
		return errors.New("Cannot move opponent's piece.")
	}

	var prevEp uint64 = game.board.ep
	game.MakeMove(move)
	game.board.ep = prevEp
	var inCheck bool = game.board.isKingInCheck(move.fromColor)
	game.UndoMove(move)
	game.board.ep = prevEp
	
	if (inCheck) {
		return errors.New("King would be in check.")
	}

	return nil
}

func (game *Game) MakeMove(move *Move) {
	switch move.flag {
	case QUIET:
		game.board.quietMove(move)
	case CAPTURE:
		game.board.capture(move)
	case KING_SIDE_CASTLE:
		game.board.castleKingSide(move)
	case QUEEN_SIDE_CASTLE:
		game.board.castleQueenSide(move)
	case PROMOTION:
		if (move.toBoard == PAWN) {
			move.toBoard = QUEEN
		}
		game.board.quietMove(move)
	case EP_CAPTURE:
		game.board.epCapture(move)
	}

	game.board.ep = move.ep
	game.board.updateCastleRights(move)
	game.turn = oppColor[game.turn]
}

func (game *Game) UndoMove(move *Move) {
	move.points *= -1
	game.MakeMove(move)
	move.points *= -1
}

func (game *Game) GetMoves() []*Move {
	var list []*Move
	var pieces uint64 = 0

	pieces = game.board.getPieces(KING, game.turn)
	list = append(list, game.getPieceMoves(pieces, game.turn, game.board.getKingSet)...)

	pieces = game.board.getPieces(QUEEN, game.turn)
	list = append(list, game.getPieceMoves(pieces, game.turn, game.board.getQueenSet)...)

	pieces = game.board.getPieces(ROOK, game.turn)
	list = append(list, game.getPieceMoves(pieces, game.turn, game.board.getRookSet)...)

	pieces = game.board.getPieces(BISHOP, game.turn)
	list = append(list, game.getPieceMoves(pieces, game.turn, game.board.getBishopSet)...)

	pieces = game.board.getPieces(KNIGHT, game.turn)
	list = append(list, game.getPieceMoves(pieces, game.turn, game.board.getKnightSet)...)

	pieces = game.board.getPieces(PAWN, game.turn)
	list = append(list, game.getPieceMoves(pieces, game.turn, game.board.getPawnSet)...)

	return list
}

func (game *Game) getPieceMoves(pieces uint64, color Color,
								getSet GetSet) []*Move {
	var list []*Move 
	for pieces != 0 {
		var piece uint64 = 1 << bitScanForward(pieces)
		var moves uint64 = getSet(piece, color)
		for moves != 0 {
			var to uint64 = 1 << bitScanForward(moves)
			move, err := game.CreateMove(piece, to)
			if err == nil {
				if move.flag == PROMOTION {
					potentialPromos := [4]Piece{QUEEN, ROOK, BISHOP, KNIGHT}
					for _, promo := range potentialPromos {
						move.toBoard = promo
						err = game.CheckMove(move)
						if err == nil {
							list = append(list, move)
						}
					}
				} else {
					err = game.CheckMove(move)
					if err == nil {
						list = append(list, move)
					}
				}
			}
			moves ^= to
		}
		pieces ^= piece
	}

	return list
}

func (game *Game) GetGameStatus() GameStatus {
	return game.status
	var moves []*Move = game.GetMoves()

	// If no legal moves, checkmate
	if (len(moves) == 0) {
		if (game.turn == WHITE) {
			return BLACK_WON
		} else {
			return WHITE_WON 
		}
	}

	return IN_PLAY
}

func (game *Game) SetGameStatus(status string) {
	switch status {
	case "0-1":
		game.status = BLACK_WON
	case "1-0":
		game.status = WHITE_WON
	case "1/2-1/2":
		game.status = DRAW
	case "*":
		game.status = IN_PLAY
	}
}