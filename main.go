package main

import (
	//"os"
	//"bufio"
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

type GameChannel interface {
	getType() string
	isUserWhite() bool
	getGameStatus() string
	getCurrMove() string
	getWinner() string
}

func main() {
	//engine := lichess.Lichess{}
	engine := DefaultEngine{}
	gameChannel := make(chan DefaultGameChannel)
	engine.Setup(gameChannel)

	var wg sync.WaitGroup

	wg.Add(2)
	go handleGame(gameChannel, &wg)
	wg.Wait()
}

func handleGame(gameChannel chan DefaultGameChannel, wg *sync.WaitGroup) {
	defer wg.Done()
	game := Game{}

	for {
		gameUpdate := <-gameChannel

		switch gameUpdate.getType() {
			case "gameFull":
				if gameUpdate.isUserWhite() {
					game.userWhite = true
					game.usersTurn = true
				} else {
					game.userWhite = false
				}
				game.board = createBoard(game.userWhite)
				printBoard(game.board, game.userWhite)
			case "gameState":
				switch gameUpdate.getGameStatus() {
					case "aborted", "resign", "timeout", "mate", "nostart":
						updateMoveList(&game, gameUpdate.getCurrMove())
						printHeader(game.numMoves)
						printBoard(game.board, game.userWhite)
						printFooter(gameUpdate.getWinner() + " wins!")
						return
					case "stalemate":
						printFooter("Stalemate!")
						return
					default:
						updateMoveList(&game, gameUpdate.getCurrMove())
						printHeader(game.numMoves)
						printBoard(game.board, game.userWhite)
						game.usersTurn = !game.usersTurn
				}
			case "chatLine":
		}

		if game.usersTurn {
			//promptAction(engine)
			//fmt.Println("\r\033[K\033[1A");
		}
		fmt.Println("\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\n")
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
