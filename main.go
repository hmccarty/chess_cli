package main

import (
	"os"
	"bufio"
	"fmt"
	"sync"
	"github.com/hmccarty/chess_cli/lichess"
)

type Engine interface {
	Setup()
	getUser() User
	getBoardChannel() <-chan GameChannel
}

type User interface {
	getID() string
}

type GameChannel interface {
	getType() string
	getMoves() string
	getWinner() string
}

func main() {
	engine := lichess.Lichess{}

	engine.authenticateClient()
}

func handleUserInput(engine Engine, gameId string, ch <-chan lichess.BoardResp, wg *sync.WaitGroup) {
	defer wg.Done()

	user := engine.getUser()
	game := Game{ID : gameId}

	for {
		boardResp := <- ch

		switch boardResp.Type {
			case "gameFull":
				if boardResp.White.ID == user.getID() {
					game.userWhite = true
					game.usersTurn = true
				} else {
					game.userWhite = false
				}
				game.board = createBoard(game.userWhite)
				printBoard(game.board, game.userWhite)
			case "gameState":
				switch boardResp.Status {
					case "aborted", "resign", "timeout", "mate", "nostart":
						updateMoveList(&game, boardResp.Moves)
						printHeader(game.numMoves)
						printBoard(game.board, game.userWhite)
						fmt.Println(boardResp.Winner)
						printFooter(boardResp.Winner + " wins!")
						return
					case "stalemate":
						printFooter("Stalemate!")
						return
					default:
						updateMoveList(&game, boardResp.Moves)
						//printHeader(game.numMoves)
						printBoard(game.board, game.userWhite)
						game.usersTurn = !game.usersTurn
				}
			case "chatLine":
		}

		if game.usersTurn {
			promptAction(game.ID)
			fmt.Println("\r\033[K\033[1A");
		}
		fmt.Println("\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\n")
	}
}

func promptAction(gameId string) {
	fmt.Print("Action (move, resign or draw): ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	path := lichessURL + gamePath + gameId + movePath + response
	path = strings.TrimSpace(path)
	_, err := client.Post(path, "plain/text", strings.NewReader(""))
	if err != nil {
		fmt.Print("Invalid option, try again\n")
		promptAction(client, gameId)
	}
}
