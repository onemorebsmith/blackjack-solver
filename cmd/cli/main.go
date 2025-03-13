package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/onemorebsmith/blackjack-solver/cmd"
	"github.com/onemorebsmith/blackjack-solver/src/blackjack/strategies"
)

type CommandLine struct {
	Decks         int     `name:"decks" default:"6"`
	H17           bool    `name:"h17" default:"false"`
	RSA           bool    `name:"rsa"`
	RoundsPerHour float32 `name:"rph" default:"100"`
	DAS           bool    `name:"das"`
	Shoes         int     `name:"shoes" default:"6"`
	Pen           float32 `name:"pen" default:"1.2"`
	Spread        string  `name:"spread"`
	Strategy      string  `name:"strat" default:"hilo"`
	Splits        int     `name:"splits" default:"3"`
	Unit          float32 `name:"unit" default:"25"`
}

func parseSpread(s string) (map[int]strategies.BidStrategy, error) {
	split := strings.Split(s, ";")
	created := map[int]strategies.BidStrategy{}
	for _, s := range split {
		parts := strings.Split(s, ":")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid spread format for key %s", s)
		}
		tc, err := strconv.ParseInt(parts[0], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed parsing %s as number", parts[0])
		}
		bid, err := strconv.ParseInt(parts[1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed parsing %s as number", parts[0])
		}
		hands := 1
		if len(parts) == 3 {
			parsedHands, err := strconv.ParseInt(parts[2], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("failed parsing %s as number", parts[0])
			}
			hands = int(parsedHands)
		}

		created[int(tc)] = strategies.BidStrategy{Hands: hands, Units: float32(bid)}
	}
	return created, nil
}

func main() {
	var commandLine CommandLine
	_ = kong.Parse(&commandLine)
	bidspread, err := parseSpread(commandLine.Spread)
	if err != nil {
		panic(err)
	}

	strategy := commandLine.Strategy
	if commandLine.Strategy == "" {
		strategy = "flatbet"
	}
	cmd.Run(cmd.BJConfig{
		Decks:         commandLine.Decks,
		IsH17:         commandLine.H17,
		IsDAS:         commandLine.DAS,
		IsRSA:         commandLine.RSA,
		ShoesToSim:    commandLine.Shoes,
		Penetration:   commandLine.Pen,
		MaxSplits:     commandLine.Splits,
		Bidspread:     bidspread,
		RoundsPerHour: commandLine.RoundsPerHour,
		Strategy:      strings.ToLower(strategy),
	})
}
