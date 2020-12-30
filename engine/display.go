package engine

import (
	"fmt"
	"github.com/fatih/color"
	"bufio"
	"os"
)

var pieceToChar = [7]string{"K", "Q", "R", "B", "N", "p", "X"}

func printHeader(moveNum int) {
	fmt.Println("~~~~~~~~~~~~~~~~~~~~")
	
	var turnLabel string
	if (moveNum % 2) == 0 {
		turnLabel = "White"
	} else {
		turnLabel = "Black"
	}

	fmt.Printf("Move #%d, %s's turn\n", moveNum, turnLabel)
}

func printFooter(result string) {
	fmt.Println("~~~~~~~~~~~~~~~~~~~~")

	fmt.Printf("%s\n", result)
}

func printRawBitBoard(board uint64) {
	for i := 8; i > 0; i-- {
	fmt.Printf("%08b\n", uint8(board >> (8 * (i - 1))))
	}
}

func printMoveList(moves *MoveList) {
	var currMove *MoveList = moves
	fmt.Printf("Moves: ")
	for {
		if (currMove == nil || currMove.next == nil) {
			break
		}
		fmt.Printf("%s, ", TranslateMove(currMove.move))
		currMove = currMove.next
	}
	fmt.Println()
}

func promptPiecePromotion() Board {
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Printf("What piece would you like to promote to? (Q, B, N, R): ")
		resp, _, _ := reader.ReadRune()
		switch resp {
		case 'Q':
			return QUEEN
		case 'R': 
			return ROOK
		case 'B':
			return BISHOP
		case 'N':
			return KNIGHT
		default:
			fmt.Println("Invalid piece, please only enter the listed options.")
		}
	}
}

func printBoard(board [7]uint64, pieceColor [2]uint64) {
		for i := uint8(8); i > 0; i-- {
			fmt.Printf("%d  ", i)
			for j := uint8(0); j < 8; j++ {
				var pos uint8 = (i * 8) - j - 1

				if (((pieceColor[WHITE] >> pos) & 1) == 1) {
					// Print white pieces if they exist on given square
					color.Set(color.FgBlue)
				} else if (((pieceColor[BLACK] >> pos) & 1) == 1) {
					// Print black pieces if they exist on given square
					color.Set(color.FgRed)
				}

				printPiece(board, pos)
				color.Unset()
			}
			color.Unset()
			fmt.Println()
		}
		fmt.Println("   A  B  C  D  E  F  G  H ")
}

func printPiece(board [7]uint64, pos uint8) {
	for pieceType, piece := range board {
		if ((piece >> pos) & 1 == 1) {
			fmt.Printf("%s  ", pieceToChar[pieceType])
			return
		}
	}
}