package main

import (
	"os"
	"bufio"
	"fmt"
	"sync"
	"strings"
	"github.com/hmccarty/gochess/goengine"
)

func main() {
	engine := goengine.GoEngine{}

	// Create games from PGN format
	//engine.scanPGN("goengine/evaluator/dataset/2017-01.bare.[7705].pgn", 1)

	// Creates new game within console
	startClientGame(engine)
}

func startClientGame(engine goengine.GoEngine) {
	inputChan := make(chan string)
	outputChan := make(chan string)
	engine.Setup(outputChan, inputChan)

	var wg sync.WaitGroup
	wg.Add(2)
	go handleGame(inputChan, outputChan, &wg)
	go engine.Run(&wg)
	wg.Wait()
}

func handleGame(inputChan chan string, outputChan chan string,
				wg *sync.WaitGroup) {
	defer wg.Done()

	reader := bufio.NewReader(os.Stdin)
	for {
		update := <-inputChan

		switch update {
			case "aborted", "resign", "timeout", "mate", "nostart":
				return
			default:
				printBoard(update)
		}

		fmt.Print("Action (move, resign or draw): ")
		response, _ := reader.ReadString('\n')
		outputChan <- strings.TrimSpace(response)
	}
}
