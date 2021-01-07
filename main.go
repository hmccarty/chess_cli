package main

import (
	"os"
	"bufio"
	"fmt"
	"sync"
	"strings"
)

func main() {
	engine := DefaultEngine{}

	// Create games from PGN format
	// engine.scanPGN("goengine/evaluator/dataset/2017-01.bare.[7705].pgn", 5000)

	// Creates new game within console
	startClientGame(engine)
}

func startClientGame(engine DefaultEngine) {
	gameChannel := make(chan DefaultGameMsg)
	inputChannel := make(chan string)
	engine.Setup(gameChannel, inputChannel)

	var wg sync.WaitGroup
	wg.Add(2)
	go handleGame(gameChannel, inputChannel, &wg)
	go engine.run(&wg)
	wg.Wait()
}

func handleGame(gameChannel chan DefaultGameMsg, inputChannel chan string,
				wg *sync.WaitGroup) {
	defer wg.Done()

	reader := bufio.NewReader(os.Stdin)
	for {
		gameUpdate := <-gameChannel

		switch gameUpdate.getType() {
			case "gameState":
				switch gameUpdate.getGameStatus() {
					case "aborted", "resign", "timeout", "mate", "nostart":
						return
				}
		}

		fmt.Print("Action (move, resign or draw): ")
		response, _ := reader.ReadString('\n')
		inputChannel <- strings.TrimSpace(response)
	}
}
