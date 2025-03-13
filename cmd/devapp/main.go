package main

import (
	"github.com/onemorebsmith/blackjack-solver/cmd"
	"github.com/onemorebsmith/blackjack-solver/src/blackjack/strategies"
)

func main() {
	cmd.Run(cmd.BJConfig{
		Decks:         6,
		IsH17:         false,
		IsRSA:         true,
		IsDAS:         true,
		Penetration:   1.5,
		RoundsPerHour: 100,
		ShoesToSim:    10000000,
		MaxSplits:     3,
		Strategy:      "hilo",
		Bidspread: map[int]strategies.BidStrategy{
			0: {Hands: 1, Units: 1},
			1: {Hands: 1, Units: 2},
			2: {Hands: 1, Units: 3},
			3: {Hands: 1, Units: 5},
			4: {Hands: 1, Units: 10},
			5: {Hands: 1, Units: 12},
		},
	})
}
