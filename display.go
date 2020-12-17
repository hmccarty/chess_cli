package main

import (
	"fmt"
	"github.com/fatih/color"
)

var pieceToChar = [6]string{"K", "Q", "R", "B", "N", "p",}

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

func printBoard(whiteBoard [6]uint64, blackBoard [6]uint64) {
		for i := uint8(8); i > 0; i-- {
			fmt.Printf("%d  ", i + 1)
			for j := uint8(0); j < 8; j++ {
				var pos uint8 = (i * 8) - j - 1

				// Print white pieces if they exist on given square
				color.Set(color.FgBlue)
				piecePrinted := printPiece(whiteBoard, pos)
				if (piecePrinted) {
					continue
				}

				// Print black pieces if they exist on given square
				color.Set(color.FgRed)
				piecePrinted = printPiece(blackBoard, pos)
				if (piecePrinted) {
					continue
				}

				// Print an empty square if nothing previously printed
				color.Unset()
				fmt.Printf("X  ")
			}
			color.Unset()
			fmt.Println()
		}
		fmt.Println("   A  B  C  D  E  F  G  H ")
}

func printPiece(board [6]uint64, pos uint8) bool {
	for pieceType, piece := range board {
		if ((piece >> pos) & 1 == 1) {
			fmt.Printf("%s  ", pieceToChar[pieceType])
			return true
		}
	}
	return false
}