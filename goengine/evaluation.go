package goengine

import "fmt"

const MAX_INT = int(^uint(0) >> 1)
const MIN_INT = -MAX_INT - 1

func minimax(game *Game, depth int, max bool,
			alpha int, beta int) (int, *Move) {
	if depth == 0 {
		var move *Move = game.moves[len(game.moves) - 1]
		return move.points, move
	}

	var best int 
	var bestMove *Move
	if max {
		best = MIN_INT
		var moves []*Move = game.getValidMoves()

		for _, move := range moves {
			game.makeMove(move)
			var value int
			value, _ = minimax(game, depth-1, false, alpha, beta)
			value += move.points
			game.undoMove()

			if value > best {
				best = value
				bestMove = move
			}

			if best > alpha {
				alpha = best
			}

			if beta <= alpha {
				break
			}
		}
		return best, bestMove
	} else {
		best = MAX_INT
		var moves []*Move = game.getValidMoves()

		for _, move := range moves {
			game.makeMove(move)
			var value int
			value, _ = minimax(game, depth-1, true, alpha, beta)
			value -= move.points
			game.undoMove()

			if value < best {
				best = value
				bestMove = move
			}

			if best < beta {
				beta = best
			}

			if beta <= alpha {
				break
			}
		}
		return best, bestMove
	}
}

func dividePerft(game *Game, depth int) int {
	if depth == 0 {
		return 1
	}

	var total int = 0
	var moves []*Move = game.getValidMoves()

	for _, move := range moves {
		game.makeMove(move)
		num := perft(game, depth - 1)
		fmt.Printf("%s: %d\n", move.ToString(), num)
		if move.ToString() == "a8a8" {
			fmt.Println(move.flag)
		}
		total += num
		game.undoMove()
	}

	return total
}

func perft(game *Game, depth int) int {
	if depth == 0 {
		return 1
	}

	var num int = 0
	var moves []*Move = game.getValidMoves()

	for _, move := range moves {
		game.makeMove(move)
		num += perft(game, depth - 1)
		game.undoMove()
	}

	return num
}