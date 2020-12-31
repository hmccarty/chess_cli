package main

import (
	"fmt"
	"github.com/fatih/color"
	//"unicode/utf8"
	//"bufio"
	//"os"
)

// func printHeader(moveNum int) {
// 	fmt.Println("~~~~~~~~~~~~~~~~~~~~")
	
// 	var turnLabel string
// 	if (moveNum % 2) == 0 {
// 		turnLabel = "White"
// 	} else {
// 		turnLabel = "Black"
// 	}

// 	fmt.Printf("Move #%d, %s's turn\n", moveNum, turnLabel)
// }

// func printFooter(result string) {
// 	fmt.Println("~~~~~~~~~~~~~~~~~~~~")

// 	fmt.Printf("%s\n", result)
// }

// func printRawBitBoard(board uint64) {
// 	for i := 8; i > 0; i-- {
// 	fmt.Printf("%08b\n", uint8(board >> (8 * (i - 1))))
// 	}
// }

// func printMoveList(moves *MoveList) {
// 	var currMove *MoveList = moves
// 	fmt.Printf("Moves: ")
// 	for {
// 		if (currMove == nil || currMove.next == nil) {
// 			break
// 		}
// 		fmt.Printf("%s, ", currMove.move.translate())
// 		currMove = currMove.next
// 	}
// 	fmt.Println()
// }

// func promptPiecePromotion() Piece {
// 	reader := bufio.NewReader(os.Stdin)
	
// 	for {
// 		fmt.Printf("What piece would you like to promote to? (Q, B, N, R): ")
// 		resp, _, _ := reader.ReadRune()
// 		switch resp {
// 		case 'Q':
// 			return QUEEN
// 		case 'R': 
// 			return ROOK
// 		case 'B':
// 			return BISHOP
// 		case 'N':
// 			return KNIGHT
// 		default:
// 			fmt.Println("Invalid piece, please only enter the listed options.")
// 		}
// 	}
// }

func printBoard(fen string) {
	var n int = 64
	fmt.Println("   A  B  C  D  E  F  G  H ")
	fmt.Printf("%d  ", (n / 8))
	for fen != "" {	
		piece := fen[0]
		fen = fen[1:]

		switch piece {
		case 'K', 'Q', 'R', 'B', 'N', 'P':
			// Print white pieces if they exist on given square
			color.Set(color.FgBlue)
			fmt.Printf("%s  ", string(piece))
			n -= 1
		case 'k', 'q', 'r', 'b', 'n', 'p':
			// Print black pieces if they exist on given square
			color.Set(color.FgRed)
			fmt.Printf("%s  ", string(piece))
			n -= 1
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			var spaces int = int(piece - 48)
			n -= spaces
			for i := 0; i < spaces; i++ {
				fmt.Printf("X  ")
			}
		case '/':
			fmt.Printf("\n%d  ", (n / 8))
		}
		color.Unset()
	}
	fmt.Println("\n   A  B  C  D  E  F  G  H ")
}

// func printBoard(board Board) {
// 		for i := uint8(8); i > 0; i-- {
// 			fmt.Printf("%d  ", i)
// 			for j := uint8(0); j < 8; j++ {
// 				var pos uint8 = (i * 8) - j - 1

// 				if (((board.color[WHITE] >> pos) & 1) == 1) {
// 					// Print white pieces if they exist on given square
// 					color.Set(color.FgBlue)
// 				} else if (((board.color[BLACK] >> pos) & 1) == 1) {
// 					// Print black pieces if they exist on given square
// 					color.Set(color.FgRed)
// 				}

// 				printPiece(board, pos)
// 				color.Unset()
// 			}
// 			color.Unset()
// 			fmt.Println()
// 		}
// 		fmt.Println("   A  B  C  D  E  F  G  H ")
// }

// func printPiece(board Board, pos uint8) {
// 	for pieceType, piece := range board.piece {
// 		if ((piece >> pos) & 1 == 1) {
// 			fmt.Printf("%s  ", pieceToChar[pieceType])
// 			return
// 		}
// 	}
// }