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

func printBoard(board [8][8]*Piece) {
		for i := 0; i < 8; i++ {
			fmt.Printf("%d  ", 8 - i)
			for j := 0; j < 8; j++ {
				if board[i][j] == nil {
					color.Unset()
				} else if board[i][j].color == WHITE {
					color.Set(color.FgBlue)
				} else if board[i][j].color == BLACK {
					color.Set(color.FgRed)
				}

				if board[i][j] == nil {
					fmt.Printf("X  ")
				} else {
					fmt.Printf("%s  ", pieceToChar[board[i][j].class])
				}
			}
			color.Unset()
			fmt.Println()
		}
		fmt.Println("   A  B  C  D  E  F  G  H ")
}