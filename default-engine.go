package main

import (
	"sync"
	"fmt"
	"github.com/hmccarty/gochess/goengine"
)

type DefaultEngine struct {
	user *DefaultUser
	game *goengine.Game
	gameChannel chan DefaultGameMsg
	inputChannel chan string
}

type DefaultUser struct {
	id string
}

type DefaultGameMsg struct {
	msgType string
	userWhite bool
	gameStatus string
	currMove string
	winner string
}

func (engine *DefaultEngine) Setup(gameChannel chan DefaultGameMsg,
								   inputChannel chan string) {

	engine.gameChannel = gameChannel
	engine.inputChannel = inputChannel

	engine.game = &goengine.Game{}
	engine.game.Setup()
}

func (engine *DefaultEngine) getUser() *DefaultUser {
	return engine.user
}

func (engine *DefaultEngine) getGameChannel() chan DefaultGameMsg {
	return engine.gameChannel
}

func (engine *DefaultEngine) scanPGN(fileName string, numGames int) {
	games := goengine.ScanGames(fileName, numGames)
	for _, game := range games {
		fmt.Println(game.GetFENString())
	}
}

func (engine *DefaultEngine) run(wg *sync.WaitGroup) {
	defer wg.Done()
	gameMsg := DefaultGameMsg{}
	gameMsg.msgType = "gameFull"
	printBoard(engine.game.GetFENString())

	for {
		printMoveList(engine.game.GetMoves())
		engine.gameChannel <- gameMsg
		cmd := <- engine.inputChannel

		gameMsg := DefaultGameMsg{}
		gameMsg.msgType = "gameState"
		gameMsg.currMove = cmd

		from, to := engine.game.ProcessCommand(cmd)
		move, err := engine.game.CreateMove(from, to)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = engine.game.CheckMove(move)
		if err != nil {
			fmt.Println(err)
			continue
		}

		engine.game.MakeMove(move)
		printBoard(engine.game.GetFENString())
		var gameStatus goengine.GameStatus = engine.game.GetGameStatus()
		switch (gameStatus) {
		case goengine.WHITE_WON:
			fmt.Println("White won!")
			gameMsg.gameStatus = "mate"
			engine.gameChannel <- gameMsg
			return
		case goengine.BLACK_WON:
			fmt.Println("Black won!")
			gameMsg.gameStatus = "mate"
			engine.gameChannel <- gameMsg
			return
		case goengine.DRAW:
			fmt.Println("Draw!")
			gameMsg.gameStatus = "mate"
			engine.gameChannel <- gameMsg
			return
		}
	}
}

func (gameMsg *DefaultGameMsg) getType() string {
	return gameMsg.msgType
}

func (gameMsg *DefaultGameMsg) isUserWhite() bool {
	return gameMsg.userWhite
}

func (gameMsg *DefaultGameMsg) getGameStatus() string {
	return gameMsg.gameStatus
}

func (gameMsg *DefaultGameMsg) getCurrMove() string {
	return gameMsg.currMove
}

func (gameMsg *DefaultGameMsg) getWinner() string {
	return gameMsg.winner
}