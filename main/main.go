package main

import (
	"log"
	"sync"

	"github.com/onemorebsmith/blackjack-solver/blackjack"
	"github.com/onemorebsmith/blackjack-solver/blackjack/core"
	"github.com/onemorebsmith/blackjack-solver/blackjack/strategies"
)

// will run iterations * handsPerGame times
const shoesPerGame = 10000
const iterations = 10

func PlayGame(rules *blackjack.BlackjackGameRules, shoes int, bankrole float32) blackjack.GameResults {
	deck := core.GenerateShoe(6).Shuffle()
	deck.PreviewCard = func(c core.Card) {
		rules.TrackingStrategy.Update(c)
	}
	totalGames := 0

	results := []blackjack.GameResults{}
	for i := 0; i < shoes; i++ {
		result := blackjack.PlayShoe(deck, rules, bankrole)
		if hl, ok := rules.TrackingStrategy.(*strategies.HighLowCountStrategy); ok {
			result.AvgTC = hl.AggregatedTC / float32(hl.Updates)
			result.HighTC = hl.HighTC
			result.LowTC = hl.LowTC
		}

		results = append(results, result)
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
	overallResults := make([]blackjack.GameResults, iterations)
	// sync := make(chan bool, 16)
	// for i := 0; i < 16; i++ {
	// 	sync <- true
	// }
	wg := sync.WaitGroup{}

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			overallResults[idx] = PlayGame(bjRules, shoesPerGame, 10000)
		}(i)
	}
	wg.Wait()

	aggregatedResults := blackjack.AggregateResults(overallResults...)
	log.Println("====================================")
	log.Printf("%d Deck, %f pen, %d hands", 6, bjRules.Penetration, aggregatedResults.Hands)
	log.Printf("   EV (units):       %f", aggregatedResults.EV)
	log.Printf("   EV (hand):        %f", aggregatedResults.EV/float32(aggregatedResults.Hands))
	log.Printf("   EV (100 hands):   %f", aggregatedResults.EV/float32(aggregatedResults.Hands)*100*50)
	log.Printf("   W/L/P:            %d/%d/%d", aggregatedResults.Wins, aggregatedResults.Losses, aggregatedResults.Pushes)
	log.Printf("   Blackjacks:       %d", aggregatedResults.Blackjacks)
	log.Printf("   Blackjack (pct):  %f", float32(aggregatedResults.Blackjacks)/float32(aggregatedResults.Hands))
	log.Printf("TC Stats --- ")
	log.Printf("   HighTC (avg)      %f ", aggregatedResults.HighTC/float32(aggregatedResults.Hands))
	log.Printf("   LowTC  (avg)      %f ", aggregatedResults.LowTC/float32(aggregatedResults.Hands))
	log.Printf("   AvgTC  (avg)      %f ", aggregatedResults.AvgTC/float32(aggregatedResults.Hands))

	log.Printf("Overall: %+v", aggregatedResults)
}
