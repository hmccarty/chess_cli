package tests

import (
	"testing"
	"github.com/hmccarty/gochess/goengine"
)

func TestPGN(t *testing.T) {
	engine := goengine.GoEngine{}
	fen := []string {
		"8/7p/2r1p3/p7/k3pNpP/1p2N3/1P2K3/8 b - - 1 60",
		"8/8/7P/6k1/p3b3/8/1K6/8 w - - 0 66",
		"8/8/8/4k2p/5P2/8/5PK1/8 b - - 0 59",
	}

	var games []string = engine.PGNToFEN("files/pgn_data.pgn", 3)
	for i, _ := range games {
		if games[i] != fen[i] {
			t.Errorf("Failed to parse PGN, got: %s, expected: %s", games[i], fen[i])
		}
	}
}