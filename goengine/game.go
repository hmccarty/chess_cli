package goengine

import (
	//"fmt"
	"errors"
)

type Game struct {
	board *Board
	turn Color
	points [2]int8
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
	var cmdData []byte = []byte(cmd)
	var startCol uint8 = 8 - (cmdData[0] - ASCII_COL_OFFSET)
	var startRow uint8 = cmdData[1] - ASCII_ROW_OFFSET
	var endCol uint8 = 8 - (cmdData[2] - ASCII_COL_OFFSET)
	var endRow uint8 = cmdData[3] - ASCII_ROW_OFFSET
	var from uint64 = 1 << ((startRow * 8) + startCol)
	var to uint64 = 1 << ((endRow * 8) + endCol)
	return from, to
}

func (game *Game) CreateMove(from uint64, to uint64) (*Move, error) {
	var move *Move = new(Move)
	move.from = new(State)
	move.from.board = from
	move.to = new(State)
	move.to.board = to
	return move, game.board.processMove(move)
}

func (game *Game) CheckMove(move *Move) error {
	if (move.from.color != game.turn) {
		return errors.New("Cannot move opponent's piece.")
	}

	game.MakeMove(move)
	var inCheck bool = game.board.isKingInCheck(move.from.color)
	game.UndoMove(move)
	
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
		if move.to.piece == PAWN {
			move.to.piece = QUEEN
		}
		game.board.quietMove(move)
	case EP_CAPTURE:
		if game.turn == WHITE {
			game.board.piece[PAWN] ^= moveNorth(move.from.enPassant)
		} else {
			game.board.piece[PAWN] ^= moveSouth(move.from.enPassant)
		}
		game.board.quietMove(move)
	}

	game.board.updateCastleRights(move)
	game.board.enPassant = move.to.enPassant
	game.points[WHITE] = move.to.points[WHITE]
	game.points[BLACK] = move.to.points[BLACK]
	game.turn = oppColor[game.turn]
}

func (game *Game) UndoMove(move *Move) {
	var tmpTo *State = move.to
	move.to = move.from
	move.from = tmpTo
	game.MakeMove(move)
	move.from = move.to
	move.to = tmpTo
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
						move.to.piece = promo
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