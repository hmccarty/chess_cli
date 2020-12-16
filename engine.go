package main

import (
	"sync"
	"fmt"
)

type DefaultEngine struct {
	user *DefaultUser
	game *Game
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

	engine.game = &Game{}
	engine.game.Setup()
}

func (engine *DefaultEngine) getUser() *DefaultUser {
	return engine.user
}

func (engine *DefaultEngine) getGameChannel() chan DefaultGameMsg {
	return engine.gameChannel
}

func (engine *DefaultEngine) run(wg *sync.WaitGroup) {
	defer wg.Done()
	gameMsg := DefaultGameMsg{}
	gameMsg.msgType = "gameFull"
	engine.gameChannel <- gameMsg

	for {
		move := <- engine.inputChannel
		fmt.Printf("Received move: %s\n", move)
		engine.game.AddNewMove(move)

		gameMsg := DefaultGameMsg{}
		gameMsg.msgType = "gameState"
		gameMsg.currMove = move
		engine.gameChannel <- gameMsg
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