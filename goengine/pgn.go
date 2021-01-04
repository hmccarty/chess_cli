package goengine

import (
	"bufio"
	"fmt"
	"regexp"
	"os"
	"io"
)

func ScanGames(fileName string, numGames int) []Game {
	var games []Game

	file, err := os.Open(fileName)
    if err != nil {
		return nil
    }
	defer file.Close()

    reader := bufio.NewReader(file)
    for i := 0; i < numGames; i++ {
		var game Game = Game{}
		game.Setup()
		var line string
	
		for {
			line, err = reader.ReadString('\n')
			if (err != nil && err != io.EOF) ||
			    line == "\r\n" {
				break
			}
		}

		if err != nil {
			break
		}

		var data string = ""
		for {
			line, err = reader.ReadString('\n')
			if (err != nil && err != io.EOF) ||
			    line == "\r\n" {
				break
			}
			data += line
		}

		re := regexp.MustCompile(`[A-Za-z]+[\w-\+]+`)
		var moves []string = re.FindAllString(data, -1)

		for _, cmd := range moves {
			from, to := game.ProcessCommand(cmd)
			move, err := game.CreateMove(from, to)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			err = game.CheckMove(move)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			game.MakeMove(move)

			if err != nil {
				fmt.Println(err)
				break
			}
		}

		games = append(games, game)
	}
	
	return games
}