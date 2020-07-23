package main

import (
	"io"
	"log"
	"fmt"
	"os"
	//"net/url"
	"encoding/json"

	"github.com/nmrshll/oauth2-noserver"
	"golang.org/x/oauth2"
)

type Event struct {
	Type string `json:"type,omitempty"`
	game Game `json:"game,omitempty"`
}

type Game struct {
	ID string `json:"id,omitempty"`
}

type Board struct {
	Type string `json:"type,omitempty"`
	Moves string `json:"moves,omitempty"`
}

const lichessURL = "https://lichess.org"
const streamEventPath = "/api/stream/event"
const seekPath = "/api/board/seek"
const streamBoardPath = "/api/board/game/stream/"

func main() {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("LICHESS_CLIENT_ID"),            // also known as client key sometimes
		ClientSecret: os.Getenv("LICHESS_CLIENT_SECRET"), // also known as secret key
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
	
	event := new(Event)
	getJSONStream(client, lichessURL + streamEventPath, event)
	fmt.Println(event.game.ID)
	board := new(Board)
	getJSONStream(client, lichessURL + streamBoardPath + (event.game.ID), board)
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
	fmt.Println("First")
	dec := json.NewDecoder(resp.Body)
	for {
		err := dec.Decode(target)
		if err != nil {
			if err == io.EOF {
		 		break
			}
			log.Fatal(err)
		}
		fmt.Println("Second")
		t, ok := target.(*Event)
		fmt.Println("Third")
		if ok {
			fmt.Println("Fourth")
			if t.Type == "gameStart" {
				break
			}
		} else {
			fmt.Println("Fourth")
			t, ok := target.(*Board)
			fmt.Println("test")
			if ok {
				if t.Type == "gameState" {
					fmt.Println(t.Moves)
				}
			}
		}
	}
}