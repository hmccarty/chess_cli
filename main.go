package main

import (
	"io"
	"bufio"
	"log"
	"fmt"
	"sync"
	"strings"
	"os"
	"encoding/json"

	"github.com/fatih/color"
	"github.com/nmrshll/oauth2-noserver"
	"golang.org/x/oauth2"
)

type User struct {
	ID string `json:"id"`
	Username string `json:"username"`
	Title string `json:"title"`
}

type EventResp struct {
	Type string `json:"type"`
	Challenge ChallengeResp `json:"challenge,omitempty"`
	Game GameResp `json:"game,omitempty"`
}

type ChallengeResp struct {
	ID string `json:"id"`
	Status string `json:"created"`
	Challenger Challenger `json:"challenger"`
	Variant Variant `json:"variant"`
}

type Challenger struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Title string `json:"title"`
	Rating int `json:"rating"`
	Patron bool `json:"patron"`
	Online bool `json:"online"`
	Lag int `json:"lag"`
}

type Variant struct {
	Key string `json:"key"`
	Name string `json:"name"`
	Short string `json:"short"`
}

type GameReq struct {
	rated bool
	time int
	increment int
	variant string
	color string
	ratingRange string
}

type GameResp struct {
	ID string `json:"id"`
}

type BoardResp struct {
	Type string `json:"type"`
	Moves string `json:"moves,omitempty"`
	Status string `json:"status,omitempty"`
	White WhiteSide `json:"white,omitempty"`
	Black BlackSide `json:"black,omitempty"`
	Winner string `json:"winner,omitempty"`
}

type WhiteSide struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

type BlackSide struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

type Game struct {
	ID string
	userWhite bool
	board [8][8]byte
	moves *Move
	numMoves int
	usersTurn bool
}

type Move struct {
	data string
	nextMove *Move
}

const lichessURL = "https://lichess.org"
const accountPath = "/api/account"
const streamEventPath = "/api/stream/event"
const seekPath = "/api/board/seek"
const streamBoardPath = "/api/board/game/stream/"
const challengePath = "/api/challenge/"
const gamePath = "/api/board/game/"
const movePath = "/move/"

const WHITE = 0x00
const BLACK = 0x80

const EMPTY = 0x00
const KING = 0x01
const QUEEN = 0x02
const ROOK = 0x03
const BISHOP = 0x04
const KNIGHT = 0x05
const PAWN = 0x06

const ASCII_ROW_OFFSET = 49
const ASCII_COL_OFFSET = 97

var pieceToChar = map[byte]string{
	EMPTY: "x",
	KING : "K",
	QUEEN: "Q",
	ROOK: "R",
	BISHOP: "B",
	KNIGHT: "N",
	PAWN: "p",
}

func main() {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("LICHESS_CLIENT_ID"),
		ClientSecret: os.Getenv("LICHESS_CLIENT_SECRET"),
		Scopes:       []string{"preference:read",
		                       "challenge:read", "challenge:write",
							   "bot:play", "board:play"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://oauth.lichess.org/oauth/authorize",
			TokenURL: "https://oauth.lichess.org/oauth",
		},
	}

	client, err := oauth2ns.AuthenticateUser(conf)
	if err != nil {
		log.Fatal(err)
	}
	event := waitForGame(client)
	ch := make(chan BoardResp)
	
	var wg sync.WaitGroup
	wg.Add(2)

	go watchForGameUpdates(client, event.Game.ID, ch, &wg)
	go handleUserInput(client, event.Game.ID, ch, &wg)

	wg.Wait()
}

func createBackRank(c byte) [8]byte {
	backRank := [8]byte {ROOK | c, KNIGHT | c, BISHOP | c, QUEEN | c, KING | c, BISHOP | c, KNIGHT | c, ROOK | c}
	return backRank
}

func createPawnRank(c byte) [8]byte {
	pawnRank := [8]byte {PAWN | c, PAWN | c, PAWN | c, PAWN | c, PAWN | c, PAWN | c, PAWN | c, PAWN | c}
	return pawnRank
}

func createEmptyRank() [8]byte {
	emptyRank := [8]byte {}
	return emptyRank
}

func createBoard(isWhite bool) [8][8]byte {
	var userColor byte
	var oppColor byte

	if isWhite {
		userColor = WHITE
		oppColor = BLACK
	} else {
		userColor = BLACK
		oppColor = WHITE
	}

	board := [8][8]byte{createBackRank(oppColor),
				        createPawnRank(oppColor),
				 		createEmptyRank(),
				 		createEmptyRank(),
				 		createEmptyRank(),
				 		createEmptyRank(),
				 		createPawnRank(userColor),
				 		createBackRank(userColor),}
	return board
}

func printHeader(moveNum int) {
	fmt.Println("~~~~~~~~~~~~~~~~~~~~")
	
	var turnLabel string
	if (moveNum % 2) == 0 {
		turnLabel = "White"
	} else {
		turnLabel = "Black"
	}

	fmt.Printf("Move #%d, %s's turn\n", moveNum, turnLabel)
}

func printFooter(result string) {
	fmt.Println("~~~~~~~~~~~~~~~~~~~~")

	fmt.Printf("%s\n", result)
}



func printBoard(board [8][8]byte, isUserWhite bool) {
	if isUserWhite {
		for i := 0; i < 8; i++ {
			fmt.Printf("%d  ", 8 - i)
			for j := 0; j < 8; j++ {
				if (board[i][j] == 0) {
					color.Unset()
				} else if (0x80 & board[i][j]) == WHITE {
					color.Set(color.FgBlue)
				} else if (0x80 & board[i][j]) == BLACK {
					color.Set(color.FgRed)
				}
				fmt.Printf("%s  ", pieceToChar[0x0F & board[i][j]])
			}
			color.Unset()
			fmt.Println()
		}
		fmt.Println("   A  B  C  D  E  F  G  H ")
	} else {
		for i := 7; i >= 0; i-- {
			fmt.Printf("%d  ", 8 - i)
			for j := 7; j >= 0; j-- {
				if (board[i][j] == 0) {
					color.Unset()
				} else if (0x80 & board[i][j] == WHITE) {
					color.Set(color.FgRed)
				} else if (0x80 & board[i][j] == BLACK) {
					color.Set(color.FgBlue)
				}
				fmt.Printf("%s  ", pieceToChar[0x0F & board[i][j]])
			}
			color.Unset()
			fmt.Println()
		}
		fmt.Println("   H  G  F  E  D  C  B  A ")
	}
}

func updateMoveList(game *Game, moves string) {
	moveArr := strings.Split(moves, " ")
	moveData := moveArr[len(moveArr) - 1]
	newMove := Move{data : moveData}
	appendMove(&(game.moves), &newMove)
	game.numMoves += 1
	completeMove(&game.board, moveData)
}

func appendMove(lastMove **Move, newMove *Move) {
	if *lastMove == nil {
		*lastMove = newMove
	} else if (*lastMove).nextMove == nil {
		appendMove(&(*lastMove).nextMove, newMove)
	}

	(*lastMove).nextMove = newMove
}

func completeMove(board *[8][8]byte, move string) {
	moveData := []byte(move)

	startCol := moveData[0] - ASCII_COL_OFFSET
	startRow := 7 - (moveData[1] - ASCII_ROW_OFFSET)

	endCol := moveData[2] - ASCII_COL_OFFSET
	endRow := 7 - (moveData[3] - ASCII_ROW_OFFSET)

	piece := (*board)[startRow][startCol]
	(*board)[startRow][startCol] = EMPTY
	(*board)[endRow][endCol] = piece
}

func getUser(client *oauth2ns.AuthorizedClient) User {
	resp, err := client.Get(lichessURL + accountPath)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	user := User{}
	err = dec.Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	return user
}

func waitForGame(client *oauth2ns.AuthorizedClient) EventResp {
	resp, err := client.Get(lichessURL + streamEventPath)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	event := EventResp{}
	for {
		err := dec.Decode(&event)
		if err != nil {
			log.Fatal(err)
		}

		switch event.Type {
			case "gameStart":
				return event
			case "challenge":
				fmt.Printf("Challenge from %s\n", event.Challenge.Challenger.Name)
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Do you accept? (y or n): ")
				response, _ := reader.ReadString('\n')

				if response == "y" {
					client.Post(lichessURL + challengePath + event.Challenge.ID + "/accept", "plain/text", strings.NewReader(""))
				} else if response == "n" {
					client.Post(lichessURL + challengePath + event.Challenge.ID + "/decline", "plain/text", strings.NewReader(""))
				} else {
					fmt.Println("Invalid response")
				}
		}
	}
}

func seekGame(client *oauth2ns.AuthorizedClient, gameReq GameReq) {

}

func watchForGameUpdates(client *oauth2ns.AuthorizedClient, gameId string, ch chan<- BoardResp, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := client.Get(lichessURL + streamBoardPath + gameId)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	for {
		boardResp := BoardResp{}
		err := dec.Decode(&boardResp)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatal(err)
		}

		ch <- boardResp
	}
}

func handleUserInput(client *oauth2ns.AuthorizedClient, gameId string, ch <-chan BoardResp, wg *sync.WaitGroup) {
	defer wg.Done()

	user := getUser(client)
	game := Game{ID : gameId}

	for {
		boardResp := <- ch

		switch boardResp.Type {
			case "gameFull":
				if boardResp.White.ID == user.ID {
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
			promptAction(client, game.ID)
			fmt.Println("\r\033[K\033[1A");
		}
		fmt.Println("\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\033[1A\n")
	}
}

func promptAction(client *oauth2ns.AuthorizedClient, gameId string) {
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
