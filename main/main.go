package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/onemorebsmith/blackjack-solver/blackjack"
	"github.com/onemorebsmith/blackjack-solver/blackjack/core"
	"github.com/onemorebsmith/blackjack-solver/blackjack/strategies"
)

// will run iterations * handsPerGame times
const shoesToSimulate = 1000000

var threads = runtime.NumCPU()
var shoesPerThread = shoesToSimulate / threads

func PlayGame(rules blackjack.BlackjackGameRules, decks int, shoes int, bankrole float32, handsPerHour float32) blackjack.GameResults {
	deck := core.GenerateShoe(decks).Shuffle()
	// create a new instance of the tracking strategy as to not share state
	// with the other threads
	rules.TrackingStrategy = rules.TrackingStrategy.Instance()

	deck.PreviewCard = func(c core.Card) {
		rules.TrackingStrategy.Update(c)
	}
	totalGames := 0

	handAVs := make([]float32, 0, shoes*50) // shoes average ~45 hands heads up
	aggregatedResults := blackjack.GameResults{}
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
		aggregatedResults = blackjack.AggregateResults(aggregatedResults, result)
		handAVs = append(handAVs, result.HandAVs...)
	}
	// calculate the population variance
	evAgg := float32(0)
	handsGroupedHourly := []float32{}

	hourlyHandCounter := float32(0)
	hourlyAgg := float32(0)
	hourlyOverallTotal := float32(0)
	for _, av := range handAVs {
		evAgg += av
		hourlyAgg += av
		hourlyHandCounter++
		if hourlyHandCounter > 100 {
			hourlyOverallTotal += hourlyAgg
			handsGroupedHourly = append(handsGroupedHourly, hourlyAgg)
			hourlyHandCounter = 0
			hourlyAgg = 0
		}
	}
	varianceAgg := float32(0)
	averageEV := evAgg / float32(len(handAVs))
	for _, r := range handAVs {
		varianceAgg += (r - averageEV) * (r - averageEV)
	}
	hourlyVariance := float32(0)
	hourlyOverallTotal /= float32(len(handsGroupedHourly))
	for _, r := range handsGroupedHourly {
		hourlyVariance += (r - averageEV) * (r - averageEV)
	}

	aggregatedResults.HourlyEVVariance = float32(math.Sqrt(float64(hourlyVariance) / float64(len(handsGroupedHourly))))
	aggregatedResults.EVVariance = float32(math.Sqrt(float64(varianceAgg) / float64(len(handAVs))))
	aggregatedResults.HandAVs = nil // save some mem
	return aggregatedResults
}

const roundsPerHour = 100
const decks = 6

func main() {
	start := time.Now()
	bjRules := blackjack.NewBlackjackGameRules(blackjack.InitGame(blackjack.H17Rules, blackjack.H17Splits))
	bjRules.SetDealerHitsSoft17(false)
	bjRules.SetDoubleAfterSplit(true)
	bjRules.SetResplitAces(true)
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
			overallResults[idx] = PlayGame(*bjRules, decks, shoesPerThread, 10000, roundsPerHour)
		}(i)
	}
	wg.Wait()

	variance := float32(0)
	hourlyVariance := float32(0)
	// extract variance before we re-aggregate
	for _, r := range overallResults {
		variance += r.EVVariance
		hourlyVariance += r.HourlyEVVariance
	}
	variance /= float32(len(overallResults))
	hourlyVariance /= float32(len(overallResults))

	game := fmt.Sprintf("%d deck ", decks)
	if bjRules.DealerHitsSoft17 {
		game += "H17 "
	} else {
		game += "S17 "
	}
	if bjRules.DoubleAfterSplit {
		game += "DAS "
	}
	if bjRules.ReSplitAces {
		game += "RSA "
	}
	game = strings.TrimRight(game, " ")

	aggregatedResults := blackjack.AggregateResults(overallResults...)
	winPct := float32(aggregatedResults.Wins) / float32(aggregatedResults.Hands)
	losePct := float32(aggregatedResults.Losses) / float32(aggregatedResults.Hands)
	pushPct := float32(aggregatedResults.Pushes) / float32(aggregatedResults.Hands)
	bjPct := float32(aggregatedResults.Blackjacks) / float32(aggregatedResults.Hands)

	log.Println("====================================")
	log.Printf("   Threads %d, elapsed: %s", threads, time.Since(start).Truncate(time.Millisecond))
	log.Println("====================================")
	log.Printf("%s, %f pen, %d hands, %d rph", game, bjRules.Penetration, aggregatedResults.Hands, roundsPerHour)
	log.Printf("   EV (units):         %f units", aggregatedResults.EV)
	log.Printf("   EV (hand):          %f units", aggregatedResults.EV/float32(aggregatedResults.Hands))
	log.Printf("   EV (hourly):        %f units", aggregatedResults.EV/float32(aggregatedResults.Hands)*roundsPerHour)
	log.Printf("   W/L/P:              %f/%f/%f", winPct, losePct, pushPct)
	log.Printf("   Blackjacks:         %d, %f%%", aggregatedResults.Blackjacks, bjPct)
	log.Printf("   1 STD (hand):     +-%f units", variance)
	log.Printf("   1 STD (hourly):   +-%f units", hourlyVariance)
	log.Printf("TC Stats --- ")
	log.Printf("   HighTC (avg)        %f ", aggregatedResults.HighTC/float32(aggregatedResults.Hands))
	log.Printf("   LowTC  (avg)        %f ", aggregatedResults.LowTC/float32(aggregatedResults.Hands))
	log.Printf("   AvgTC  (avg)        %f ", aggregatedResults.AvgTC/float32(aggregatedResults.Hands))

	// dumpTCFreqTable(aggregatedResults)
}

func dumpTCFreqTable(results blackjack.GameResults) {
	type bidKv struct {
		TC   int
		Freq int
	}
	bids := []bidKv{}
	totalBids := 0
	for k, v := range results.BidsByTC {
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
}
