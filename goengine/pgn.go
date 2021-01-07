package goengine

import (
	"bufio"
	"fmt"
	"regexp"
	"os"
	"io"
)

func scanGames(fileName string, numGames int) []Game {
	var games []Game

	file, err := os.Open(fileName)
    if err != nil {
		return nil
    }
	defer file.Close()

	reader := bufio.NewReader(file)
	var line string
    for i := 0; i < numGames; i++ {
		// PGN Header
		var header map[string]string = make(map[string]string)
		keyRE := regexp.MustCompile(`\[(\w)+`)
		valueRE := regexp.MustCompile(`"[^"]+`)
		for {
			line, err = reader.ReadString('\n')
			if (err != nil && err != io.EOF) ||
			    line == "\r\n" {
				break
			}
			header[keyRE.FindString(line)[1:]] = valueRE.FindString(line)[1:]
		}

		if err != nil {
			break
		}

		// PGN Move Data
		var data string = ""
		for {
			line, err = reader.ReadString('\n')
			if (err != nil && err != io.EOF) ||
			    line == "\r\n" {
				break
			}
			data += line
		}

		// Setup Game struct
		var game Game = Game{}
		game.setup()
		game.setGameStatus(header["Result"])

		// Get list of moves in an algebraic format
		moveRE := regexp.MustCompile(`[A-Za-z][\w-]+[\+]?`)
		var moves []string = moveRE.FindAllString(data, -1)

		// Add each move into create Game struct
		var err error
		for _, cmd := range moves {
			err = game.pushSAN(cmd)
			if err != nil {
				fmt.Println(cmd)
				fmt.Println(err)
				break
			}
		}

		games = append(games, game)
	}
	
	return games
}