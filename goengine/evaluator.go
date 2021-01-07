package goengine

func alphaBetaSearch(game *Game) *Move {
	var alpha int = 0
	var best *Move = nil
	var moves []*Move = game.getValidMoves()
	for _, move := range moves {
		alpha += move.points
		game.makeMove(move)
		var respMoves []*Move = game.getValidMoves()
		for _, respMove := range respMoves {
			
		}
	}
} 