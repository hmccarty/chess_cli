package main

import (
	"fmt"
	"github.com/fatih/color"
)

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

func printBoard(board [8][8]byte, isUserWhite bool) {
	if isUserWhite {
		for i := 0; i < 8; i++ {
			fmt.Printf("%d  ", 8 - i)
			for j := 0; j < 8; j++ {
				if (board[i][j] == 0) {
					color.Unset()
				} else if (0x80 & board[i][j]) == WHITE {
					color.Set(color.FgBlue)
				} else if (0x80 & board[i][j]) == BLACK {
					color.Set(color.FgRed)
				}
				fmt.Printf("%s  ", pieceToChar[0x0F & board[i][j]])
			}
			color.Unset()
			fmt.Println()
		}
		fmt.Println("   A  B  C  D  E  F  G  H ")
	} else {
		for i := 7; i >= 0; i-- {
			fmt.Printf("%d  ", 8 - i)
			for j := 7; j >= 0; j-- {
				if (board[i][j] == 0) {
					color.Unset()
				} else if (0x80 & board[i][j] == WHITE) {
					color.Set(color.FgRed)
				} else if (0x80 & board[i][j] == BLACK) {
					color.Set(color.FgBlue)
				}
				fmt.Printf("%s  ", pieceToChar[0x0F & board[i][j]])
			}
			color.Unset()
			fmt.Println()
		}
		fmt.Println("   H  G  F  E  D  C  B  A ")
	}
}