package main

import (
	"log"

	"github.com/onemorebsmith/blackjack-solver/blackjack"
	"github.com/onemorebsmith/blackjack-solver/blackjack/core"
	"github.com/onemorebsmith/blackjack-solver/blackjack/strategies"
)

const startingBankrole = float32(1000) // in bet units
// will run iterations * handsPerGame times
const handsPerGame = 10000
const iterations = 100

func PlayGame(rules *blackjack.BlackjackGameRules, hands int64, bankrole float32) blackjack.GameResults {
	deck := core.GenerateShoe(6).Shuffle()

	totalGames := 0

	results := []blackjack.GameResults{}
	for i := 0; i < handsPerGame; i++ {
		results = append(results, blackjack.PlayShoe(deck, rules, bankrole))
		rules.TrackingStrategy.Shuffle()
		deck.Shuffle()
		totalGames++
	}

	return blackjack.AggregateResults(results...)
}

func main() {
	bjRules := blackjack.NewBlackjackGameRules(blackjack.InitGame(blackjack.H17Rules, blackjack.H17Splits))
	bjRules.SetDealerHitsSoft17(true)
	bjRules.SetDoubleAfterSplit(true) // not implemented
	bjRules.SetMaxPlayerSplits(4)
	bjRules.SetUseSimpleDeviations(false)
	bjRules.SetPenetration(1.5)
	//bjRules.TrackingStrategy = strategies.InitFlatbetStrategy()
	bjRules.TrackingStrategy = strategies.InitHighLow(map[int]strategies.BidStrategy{
		0: {Hands: 1, Units: 1},
		1: {Hands: 1, Units: 2},
		2: {Hands: 1, Units: 4},
		3: {Hands: 1, Units: 8},
		4: {Hands: 1, Units: 12},
	})
	// resultsChannel := make(chan blackjack.GameResults)
	overallResults := make([]blackjack.GameResults, 0, iterations)
	// sync := make(chan bool, 16)
	// for i := 0; i < 16; i++ {
	// 	sync <- true
	// }

	for i := 0; i < 10; i++ {
		overallResults = append(overallResults, PlayGame(bjRules, handsPerGame, startingBankrole))
	}

	aggregatedResults := blackjack.AggregateResults(overallResults...)
	log.Println("====================================")
	log.Printf("%d Deck, %f pen, %d hands", 6, bjRules.Penetration, aggregatedResults.Hands)
	log.Printf("EV (units):       %f", aggregatedResults.EV)
	log.Printf("EV (hand):        %f", aggregatedResults.EV/float32(aggregatedResults.Hands))
	log.Printf("EV (100 hands):   %f", aggregatedResults.EV/float32(aggregatedResults.Hands)*100*50)
	log.Printf("W/L/P:            %d/%d/%d", aggregatedResults.Wins, aggregatedResults.Losses, aggregatedResults.Pushes)
	log.Printf("Blackjacks:       %d", aggregatedResults.Blackjacks)
	log.Printf("Blackjack (pct):  %f", float32(aggregatedResults.Blackjacks)/float32(aggregatedResults.Hands))

	log.Printf("Overall: %+v", aggregatedResults)
}
