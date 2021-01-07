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

func (engine *GoEngine) ScanPGN(fileName string, numGames int) {
	games := scanGames(fileName, numGames)
	for _, game := range games {
		fmt.Println(game.getFENString())
	}
}

func (engine *GoEngine) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		engine.outputChan <- engine.game.getFENString()
		cmd := <- engine.inputChan

		err := engine.game.pushSAN(cmd)
		if err != nil {
			fmt.Println(err)
			continue
		}

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