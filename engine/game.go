package engine

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

type MoveList struct {
	num uint8
	move *Move
	next *MoveList
	prev *MoveList
}

func (game *Game) Setup() {
	game.board = new(Board)
	game.board.setup()
	game.turn = WHITE
}

func (game *Game) GetFENString() string {
	return game.board.getFENBoard()
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

func (game *Game) ProcessMove(fromSqr uint8, toSqr uint8) (*Move, error) {
	move, err := game.board.createMove(fromSqr, toSqr)

	if (err != nil) {
		return nil, err
	} else if (move.fromColor != game.turn) {
		return nil, errors.New("Cannot move opponent's piece.")
	}

	game.MakeMove(move)
	var inCheck bool = game.board.isKingInCheck(move.fromColor)
	game.UndoMove(move)
	
	if (inCheck) {
		return nil, errors.New("King would be in check.")
	}

	return move, nil
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
		// if (move.toBoard == PAWN) {
		// 	move.toBoard = promptPiecePromotion()
		// }
		game.board.quietMove(move)
	}

	game.board.updateCastleRights(move)
	game.turn = oppColor[game.turn]
}

func (game *Game) UndoMove(move *Move) {
	move.points *= -1
	game.MakeMove(move)
	move.points *= -1
}

func (game *Game) GetAllLegalMoves() *MoveList {
	var first *MoveList = new(MoveList)
	first.num = 1

	var pieces uint64 = 0
	var currMove *MoveList = first

	pieces = game.board.getPieces(KING, game.turn)
	currMove = game.AddMovesToList(pieces, game.turn, game.board.getKingMoves, currMove)

	pieces = game.board.getPieces(QUEEN, game.turn)
	currMove = game.AddMovesToList(pieces, game.turn, game.board.getQueenMoves, currMove)

	pieces = game.board.getPieces(ROOK, game.turn)
	currMove = game.AddMovesToList(pieces, game.turn, game.board.getRookMoves, currMove)

	pieces = game.board.getPieces(BISHOP, game.turn)
	currMove = game.AddMovesToList(pieces, game.turn, game.board.getBishopMoves, currMove)

	pieces = game.board.getPieces(KNIGHT, game.turn)
	currMove = game.AddMovesToList(pieces, game.turn, game.board.getKnightMoves, currMove)

	pieces = game.board.getPieces(PAWN, game.turn)
	game.AddMovesToList(pieces, game.turn, game.board.getPawnMoves, currMove)

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
				if (move.flag == PROMOTION) {
					
				}
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