package goengine

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

type Color uint8
const (
	WHITE Color = iota
	BLACK
)

var pieceToPoints = map[Piece]int {
	KING   : MAX_INT,
	QUEEN  : 9,
	ROOK   : 5,
	BISHOP : 3,
	KNIGHT : 3,
	PAWN   : 1,
	EMPTY  : 0,
}

var runeToPiece = map[rune]Piece {
	'K' : KING,
	'Q' : QUEEN,
	'R' : ROOK,
	'B' : BISHOP,
	'N' : KNIGHT,
}

var pieceToString = [2][7]string{{"K", "Q", "R", "B", "N", "P", "X"},
								 {"k", "q", "r", "b", "n", "p", "X"},}

var colorToString = [2]string{"w", "b"}

var oppColor = map[Color]Color {
	WHITE : BLACK,
	BLACK : WHITE,
}