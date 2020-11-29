package main

import (
	"sync"
)

type DefaultEngine struct {
	user *DefaultUser
	gameChannel chan DefaultGameChannel
}

type DefaultUser struct {
	id string
}

type DefaultGameChannel struct {
	msgType string
	userWhite bool
	gameStatus string
	currMove string
	winner string
}

func (engine *DefaultEngine) Setup(gameChannel chan DefaultGameChannel) {
	engine.gameChannel = gameChannel
}

func (engine *DefaultEngine) getUser() *DefaultUser {
	return engine.user
}

func (engine *DefaultEngine) getGameChannel() chan DefaultGameChannel {
	return engine.gameChannel
}

func (engine *DefaultEngine) run(wg *sync.WaitGroup) {
	defer wg.Done()
}

func (gameChannel *DefaultGameChannel) getType() string {
	return gameChannel.msgType
}

func (gameChannel *DefaultGameChannel) isUserWhite() bool {
	return gameChannel.userWhite
}

func (gameChannel *DefaultGameChannel) getGameStatus() string {
	return gameChannel.gameStatus
}

func (gameChannel *DefaultGameChannel) getCurrMove() string {
	return gameChannel.currMove
}

func (gameChannel *DefaultGameChannel) getWinner() string {
	return gameChannel.winner
}