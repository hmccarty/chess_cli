package goengine

import (
	"fmt"
	"errors"
)

type Game struct {
	board *Board
	moves []*Move
	turn Color
	halfmove uint8
	fullmove uint8
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
	game.fullmove = 1
	game.turn = WHITE
}

func (game *Game) GetFENString() string {
	var fen string = game.board.getFENBoard()
	fen += " "
	fen += colorToString[game.turn]
	fen += " "
	if (game.board.castle[WHITE] | game.board.castle[BLACK]) != 0 {
		if (game.board.castle[WHITE] & KING_CASTLE_MASK) != 0 {
			fen += "K"
		}
		if (game.board.castle[WHITE] & QUEEN_CASTLE_MASK) != 0 {
			fen += "Q"
		}
		if (game.board.castle[BLACK] & QUEEN_CASTLE_MASK) != 0 {
			fen += "k"
		}
		if (game.board.castle[BLACK] & QUEEN_CASTLE_MASK) != 0 {
			fen += "q"
		}
	} else {
		fen += "-"
	}
	fen += " "
	if game.board.ep != 0 {
		var sqr uint64 = (game.board.ep ^ game.board.piece[PAWN]) & game.board.ep
		var idx uint8 = bitScanForward(sqr)
		// Convert row and column into algebraic notation
		fen += string((idx / 8) + ASCII_COL_OFFSET)
		fen += string((idx % 8) + ASCII_ROW_OFFSET)
	} else {
		fen += "-"
	}
	fen += " "
	fen += fmt.Sprintf("%d %d", game.halfmove, game.fullmove)
	return fen
}

func (game *Game) PushSAN(cmd string) error {
	var move *Move = new(Move)

	if cmd == "O-O" {
		move.flag = KING_SIDE_CASTLE
		move.from = game.board.piece[KING] & game.board.color[game.turn]
		move.to = move.from >> 2
		err := game.HandleMove(move)
		if err != nil {
			return err
		}
		return nil
	} else if cmd == "O-O-O" {
		move.flag = QUEEN_SIDE_CASTLE
		move.from = game.board.piece[KING] & game.board.color[game.turn]
		move.to = move.from << 2
		err := game.HandleMove(move)
		if err != nil {
			return err
		}
		return nil
	}

	var cmdData []byte = []byte(cmd)

	// If move causes check, remove 
	if cmdData[len(cmdData) - 1:][0] == byte('+') {
		cmdData = cmdData[:len(cmdData) - 1]
	}

	// Find square the piece is moving to
	var toCol uint8 = 8 - (cmdData[len(cmdData) - 2:len(cmdData) - 1][0] - ASCII_COL_OFFSET)
	var toRow uint8 = byte(cmdData[len(cmdData) - 1:][0]) - ASCII_ROW_OFFSET
	move.to = 1 << ((toRow * 8) + toCol)

	// Store additional info about piece
	var additionalInfo []byte

	// Determine type of piece being moved
	var symbolExists bool
	move.fromBoard, symbolExists = runeToPiece[rune(cmdData[0])]
	if !symbolExists {
		move.fromBoard = PAWN
		additionalInfo = cmdData[0:len(cmdData) - 2]
	} else {
		additionalInfo = cmdData[1:len(cmdData) - 2]
	}

	// Search all possible pieces
	var pieces uint64 = game.board.piece[move.fromBoard] & game.board.color[game.turn]
	for pieces > 0 {
		var idx uint8 = bitScanForward(pieces)
		var piece uint64 = 1 << idx
		move.from = piece
		// If square piece is moving to is in piece range
		if (game.board.getPieceSet(move.fromBoard, piece, game.turn) & move.to) != 0 {
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
				err := game.HandleMove(move)
				if err != nil {
					return err
				}
				return nil
			}
		}

		pieces ^= piece
	}

	return errors.New("Couldn't find piece to carry out move.")
}

func (game *Game) HandleMove(move *Move) error {
	move.fullmove = game.fullmove
	move.halfmove = game.halfmove + 1
	
	err := game.board.processMove(move)
	if err != nil {
		return err
	}

	if move.fromColor != game.turn {
		return errors.New("Cannot move opponent's piece.")
	}

	game.MakeMove(move)
	if (game.board.isKingInCheck(move.fromColor)) {
		game.UndoMove()
		return errors.New("King would be in check.")
	}

	return nil
}

func (game *Game) MakeMove(move *Move) {
	// Modifying board based on move data
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
	game.board.castle[WHITE] = move.castle[WHITE]
	game.board.castle[BLACK] = move.castle[BLACK]
	game.turn = oppColor[game.turn]
	game.halfmove = move.halfmove
	game.fullmove = move.fullmove
	if game.turn == WHITE {
		game.fullmove += 1
	}
	game.moves = append(game.moves, move)
}

func (game *Game) UndoMove() {
	var move *Move = game.moves[len(game.moves) - 1]

	// Applying same board data to reverse last move
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

	// Popping last move from stack
	game.moves = game.moves[:len(game.moves) - 1]

	// Setting relevant game variables to new, last move
	if len(game.moves) > 0 {
		move = game.moves[len(game.moves) - 1]
		game.board.ep = move.ep
		game.board.castle[WHITE] = move.castle[WHITE]
		game.board.castle[BLACK] = move.castle[BLACK]
		game.halfmove = move.halfmove
		game.fullmove = move.fullmove
	} else {
		game.board.ep = 0
		game.board.castle[WHITE] = (KING_CASTLE_MASK | QUEEN_CASTLE_MASK)
		game.board.castle[BLACK] = (KING_CASTLE_MASK | QUEEN_CASTLE_MASK)
		game.halfmove = 0
		game.fullmove = 1
	}
	game.turn = oppColor[game.turn]
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
		// Get first piece set
		var piece uint64 = 1 << bitScanForward(pieces)
		var moves uint64 = getSet(piece, color)

		// Loop for every possible move in set
		for moves != 0 {
			var move *Move = new(Move)
			move.from = piece
			move.to = 1 << bitScanForward(moves)
			err := game.HandleMove(move)
			if err == nil {
				// If promotion, handle all possible promotions
				// Else undo handled move and append to list
				if move.flag == PROMOTION {
					game.UndoMove()
					potentialPromos := [3]Piece{ROOK, BISHOP, KNIGHT}
					for _, promo := range potentialPromos {
						move.toBoard = promo
						err = game.HandleMove(move)
						if err == nil {
							game.UndoMove()
							list = append(list, move)
						}
					}
				} else {
					game.UndoMove()
					list = append(list, move)
				}
			}
			moves ^= move.to
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