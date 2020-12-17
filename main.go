package main

import (
	"os"
	"bufio"
	"fmt"
	"sync"
	//"github.com/hmccarty/lichess"
)

type Engine interface {
	Setup()
	getUser()
	getGameChannel()
}

type User interface {
	getID() string
}

type GameMsg interface {
	getType() string
	isUserWhite() bool
	getGameStatus() string
	getCurrMove() string
	getWinner() string
}

func main() {
	//engine := lichess.Lichess{}
	engine := DefaultEngine{}
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
	game := Game{}
	// userColor := WHITE

	for {
		gameUpdate := <-gameChannel

		switch gameUpdate.getType() {
			case "gameFull":
				game.Setup()
				printBoard(game.whiteBoard, game.blackBoard)
			// case "gameState":
			// 	switch gameUpdate.getGameStatus() {
			// 		case "aborted", "resign", "timeout", "mate", "nostart":
			// 			game.AddNewMove(gameUpdate.getCurrMove())
			// 			printHeader(game.numMoves)
			// 			printBoard(game.board)
			// 			printFooter(gameUpdate.getWinner() + " wins!")
			// 			return
			// 		case "stalemate":
			// 			printFooter("Stalemate!")
			// 			return
			// 		default:
			// 			game.AddNewMove(gameUpdate.getCurrMove())
			// 			printHeader(game.numMoves)
			// 			printBoard(game.board)
			//	}
			case "chatLine":
		}

		fmt.Print("Action (move, resign or draw): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		inputChannel <- response
		// if game.turnColor == userColor {
			//promptAction(engine)
			//fmt.Println("\r\033[K\033[1A");
		// }
		//fmt.Println("\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\n")
	}
}

// func promptAction(engine Engine) {
// 	fmt.Print("Action (move, resign or draw): ")
// 	reader := bufio.NewReader(os.Stdin)
// 	response, _ := reader.ReadString('\n')
// 	path := lichessURL + gamePath + gameId + movePath + response
// 	path = strings.TrimSpace(path)
// 	_, err := client.Post(path, "plain/text", strings.NewReader(""))
// 	if err != nil {
// 		fmt.Print("Invalid option, try again\n")
// 		promptAction(client, gameId)
// 	}
// }
