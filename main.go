package main

import (
	"io"
	"log"
	"fmt"
	"strings"
	// "os"
	"encoding/json"

	"github.com/fatih/color"
	"github.com/nmrshll/oauth2-noserver"
	// "golang.org/x/oauth2"
)

type EventResp struct {
	Type string `json:"type"`
	Game GameResp `json:"game,omitempty"`
}

type GameResp struct {
	ID string `json:"id"`
}

type BoardResp struct {
	Type string `json:"type"`
	Moves string `json:"moves,omitempty"`
	Status string `json:"resign,omitempty"`
}

type Game struct {
	userWhite bool
	board [8][8]byte
	moves *Move
	numMoves int
}

type Move struct {
	data string
	nextMove *Move
}

const lichessURL = "https://lichess.org"
const streamEventPath = "/api/stream/event"
const seekPath = "/api/board/seek"
const streamBoardPath = "/api/board/game/stream/"

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
	
	event := new(EventResp)
	getJSONStream(client, lichessURL + streamEventPath, event)
	board := new(BoardResp)
	getJSONStream(client, lichessURL + streamBoardPath + (event.Game.ID), board)
	game := Game{userWhite : true}
	game.board = createBoard(game.userWhite)
	printBoard(game.board)
	updateMoveList(&game, "e2e4 c7c5 f2f4 d7d6 g1f3 b8c6 f1c4 g8f6 d2d3 g7g6 e1g1 f8g7")
	printBoard(game.board)
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

func printBoard(board [8][8]byte) {
	for i := 0; i < 8; i++ {
		fmt.Printf("%d  ", 8 - i)
		for j := 0; j < 8; j++ {
			if (board[i][j] == 0) {
				color.Unset()
			} else if (0xF0 & board[i][j] == WHITE) {
				color.Set(color.FgBlue)
			} else if (0xF0 & board[i][j] == BLACK) {
				color.Set(color.FgRed)
			}
			fmt.Printf("%s  ", pieceToChar[0x0F & board[i][j]])
		}
		color.Unset()
		fmt.Println()
	}
	fmt.Println("   A  B  C  D  E  F  G  H ")
}

func updateMoveList(game *Game, moves string) {
	moveArr := strings.Split(moves, " ")
	currMove := game.moves

	for i, move := range moveArr {
		if i >= game.numMoves {
			moveStruct := Move{data : move}
			appendMove(&game.moves, &moveStruct)
			game.numMoves += 1
			completeMove(&game.board, move)
		} else if move != currMove.data {
			fmt.Println("Moves corrupted, exiting")
			return
		}
	}
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

func getJSON(client *oauth2ns.AuthorizedClient, url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func getJSONStream(client *oauth2ns.AuthorizedClient, url string, target interface{}) {
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	for {
		err := dec.Decode(target)
		if err != nil {
			if err == io.EOF {
		 		break
			}
			log.Fatal(err)
		}
		t, ok := target.(*EventResp)
		if ok {
			if t.Type == "gameStart" {
				break
			}
		} else {
			t, ok := target.(*BoardResp)
			if ok {
				if t.Type == "gameState" {
					
				}
			}
		}
	}
}