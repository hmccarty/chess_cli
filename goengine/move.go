package goengine

type Move struct {
	flag Flag
	from *State
	to *State
}

func (move *Move) copy() *Move {
	var new *Move = new(Move)
	new.flag = move.flag
	new.from = move.from.copy()
	new.to = move.to.copy()
	return new
}

func (move *Move) ToString() string {
	var fromSqr uint8 = bitScanForward(move.from.board)
	var toSqr uint8 = bitScanForward(move.to.board)
	var startRow uint8 = fromSqr / 8
	var startCol uint8 = fromSqr % 8
	var endRow uint8 = toSqr / 8
	var endCol uint8 = toSqr % 8
	return (string((8 - startCol) + ASCII_COL_OFFSET) +
	        string(startRow + ASCII_ROW_OFFSET) + 
			string((8 - endCol) + ASCII_COL_OFFSET) +
			string(endRow + ASCII_ROW_OFFSET))
}

type State struct {
	board uint64
	piece Piece
	color Color
	kingCastle [2]bool
	queenCastle [2]bool	
	enPassant uint64
	points [2]int8
}

func (state *State) copy() *State {
	var new *State = new(State)
	new.board = state.board
	new.color = state.color
	new.kingCastle[WHITE] = state.kingCastle[WHITE]
	new.kingCastle[BLACK] = state.kingCastle[BLACK]
	new.queenCastle[WHITE] = state.queenCastle[WHITE]
	new.queenCastle[BLACK] = state.queenCastle[BLACK]
	new.enPassant = state.enPassant
	new.points[WHITE] = state.points[WHITE]
	new.points[BLACK] = state.points[BLACK]
	return new
}

type Flag uint8
const (
	UNKNOWN Flag = iota
	QUIET
	CAPTURE
	KING_SIDE_CASTLE
	QUEEN_SIDE_CASTLE
	EP_CAPTURE
	PROMOTION
)