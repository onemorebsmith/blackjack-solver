package main

import (
	"log"
	"runtime"
	"sort"
	"sync"

	"github.com/onemorebsmith/blackjack-solver/blackjack"
	"github.com/onemorebsmith/blackjack-solver/blackjack/core"
	"github.com/onemorebsmith/blackjack-solver/blackjack/strategies"
)

// will run iterations * handsPerGame times
const shoesToSimulate = 100000000

var threads = runtime.NumCPU()
var shoesPerThread = shoesToSimulate / threads

func PlayGame(rules blackjack.BlackjackGameRules, shoes int, bankrole float32) blackjack.GameResults {
	deck := core.GenerateShoe(6).Shuffle()
	// create a new instance of the tracking strategy as to not share state
	// with the other threads
	rules.TrackingStrategy = rules.TrackingStrategy.Instance()

	deck.PreviewCard = func(c core.Card) {
		rules.TrackingStrategy.Update(c)
	}
	totalGames := 0

	results := []blackjack.GameResults{}
	for i := 0; i < shoes; i++ {
		result := blackjack.PlayShoe(deck, &rules, bankrole)
		if hl, ok := rules.TrackingStrategy.(*strategies.HighLowCountStrategy); ok {
			result.BidsByTC = hl.BidsByTC
			result.AvgTC = hl.AggregatedTC / float32(hl.Updates)
			result.HighTC = hl.HighTC
			result.LowTC = hl.LowTC
		}
		rules.TrackingStrategy.Shuffle()
		deck.Shuffle()
		totalGames++
		results = append(results, result)
	}
	return blackjack.AggregateResults(results...)
}

func main() {
	bjRules := blackjack.NewBlackjackGameRules(blackjack.InitGame(blackjack.H17Rules, blackjack.H17Splits))
	bjRules.SetDealerHitsSoft17(true)
	bjRules.SetDoubleAfterSplit(true) // not implemented
	bjRules.SetMaxPlayerSplits(4)
	bjRules.SetUseSimpleDeviations(false)
	bjRules.SetPenetration(.5)
	bjRules.TrackingStrategy = strategies.InitHighLow(map[int]strategies.BidStrategy{
		0: {Hands: 1, Units: 1},
		1: {Hands: 1, Units: 2},
		2: {Hands: 1, Units: 3},
		3: {Hands: 1, Units: 5},
		4: {Hands: 1, Units: 10},
		5: {Hands: 1, Units: 12},
	})
	//bjRules.TrackingStrategy = strategies.InitFlatbetStrategy()
	// resultsChannel := make(chan blackjack.GameResults)
	overallResults := make([]blackjack.GameResults, threads)
	// sync := make(chan bool, 16)
	// for i := 0; i < 16; i++ {
	// 	sync <- true
	// }
	wg := sync.WaitGroup{}

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			overallResults[idx] = PlayGame(*bjRules, shoesPerThread, 10000)
		}(i)
	}
	wg.Wait()

	variance := float32(0)
	// extract variance before we re-aggregate
	for _, r := range overallResults {
		variance += r.EVVariance
	}
	variance /= float32(len(overallResults))

	aggregatedResults := blackjack.AggregateResults(overallResults...)
	log.Println("====================================")
	log.Printf("%d Deck, %f pen, %d shoes, %d hands", 6, bjRules.Penetration, shoesToSimulate, aggregatedResults.Hands)
	log.Printf("   EV (units):       %f", aggregatedResults.EV)
	log.Printf("   EV (hand):        %f", aggregatedResults.EV/float32(aggregatedResults.Hands))
	log.Printf("   EV (100 hands):   %f", aggregatedResults.EV/float32(aggregatedResults.Hands)*100*50)
	log.Printf("   W/L/P:            %d/%d/%d", aggregatedResults.Wins, aggregatedResults.Losses, aggregatedResults.Pushes)
	log.Printf("   Blackjacks:       %d", aggregatedResults.Blackjacks)
	log.Printf("   Blackjack (pct):  %f", float32(aggregatedResults.Blackjacks)/float32(aggregatedResults.Hands))
	log.Printf("   1 STD($):       +-%f", variance*50)
	log.Printf("TC Stats --- ")
	log.Printf("   HighTC (avg)      %f ", aggregatedResults.HighTC/float32(aggregatedResults.Hands))
	log.Printf("   LowTC  (avg)      %f ", aggregatedResults.LowTC/float32(aggregatedResults.Hands))
	log.Printf("   AvgTC  (avg)      %f ", aggregatedResults.AvgTC/float32(aggregatedResults.Hands))

	type bidKv struct {
		TC   int
		Freq int
	}
	bids := []bidKv{}
	totalBids := 0
	for k, v := range aggregatedResults.BidsByTC {
		totalBids += v
		bids = append(bids, bidKv{TC: k, Freq: v})
	}
	sort.Slice(bids, func(i, j int) bool {
		return bids[i].TC < bids[j].TC
	})
	log.Println("    TC |    %   | Freq")
	for _, b := range bids {
		log.Printf("   %d | %f | %d", b.TC, float32(b.Freq)/float32(totalBids), b.Freq)
	}

	//log.Printf("Overall: %+v", aggregatedResults)
}
