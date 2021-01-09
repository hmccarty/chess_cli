package goengine

type Move struct {
	flag Flag
	from uint64
	to uint64
	piece Piece
	target Piece
	color Color
	castle [2]uint8
	ep uint64
	fullmove uint8
	halfmove uint8
	points int
}

func (move *Move) copy() *Move {
	return &Move{
		flag     : move.flag,
		from     : move.from,
		to       : move.to,
		piece    : move.piece,
		target   : move.target,
		color    : move.color,
		castle   : move.castle,
		ep       : move.ep,
		halfmove : move.halfmove,
		fullmove : move.fullmove,
		points   : move.points,
	}
}

func (move *Move) ToString() string {
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

type Flag uint8
const (
	UNKNOWN Flag = iota
	QUIET
	CAPTURE
	K_CASTLE
	Q_CASTLE
	EP_CAPTURE
	PROMOTION
)