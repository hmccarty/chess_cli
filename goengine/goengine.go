package goengine

import (
	"sync"
	"fmt"
)

type GoEngine struct {
	game *Game
	inputChan chan string
	outputChan chan string
}

func (engine *GoEngine) Setup(inputChan chan string, outputChan chan string) {
	engine.outputChan = outputChan
	engine.inputChan = inputChan

	engine.game = &Game{}
	engine.game.setup()
}

func (engine *GoEngine) PGNToFEN(fileName string, numGames int) []string {
	fen := make([]string, numGames)
	games := scanGames(fileName, numGames)
	for i, game := range games {
		fen[i] = game.getFENString()
	}
	return fen
}

func (engine *GoEngine) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	engine.game.setFENString("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - -")

	for {
		//if engine.game.turn == WHITE {
			fmt.Println(dividePerft(engine.game, 6))
			fmt.Println(engine.game.getFENString())
			engine.outputChan <- "client " + engine.game.getFENString()
			cmd := <- engine.inputChan

			err := engine.game.pushSAN(cmd)
			if err != nil {
				fmt.Println(err)
				continue
			}
		//} else {
			// engine.outputChan <- "engine " + engine.game.getFENString()
			// _, move := minimax(engine.game, 5, true, MIN_INT, MAX_INT)
			// engine.game.makeMove(move)
		//}
		var gameStatus GameStatus = engine.game.getGameStatus()
		switch (gameStatus) {
		case WHITE_WON:
			fmt.Println("White won!")
			engine.outputChan <- "mate"
			return
		case BLACK_WON:
			fmt.Println("Black won!")
			engine.outputChan <- "mate"
			return
		case DRAW:
			fmt.Println("Draw!")
			engine.outputChan <- "mate"
			return
		}
	}
}