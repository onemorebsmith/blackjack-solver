package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/onemorebsmith/blackjack-solver/blackjack"
)

const startingBankrole = float32(1000) // in bet units
// will run iterations * handsPerGame times
const handsPerGame = 10000
const iterations = 100

type GameResults struct {
	Hands      int
	WinRate    float32
	LoseRate   float32
	PushRate   float32
	Blackjacks float32
	EV         float32
	Result     float32
	AvgTc      float32
}

func PlayGame(rules *blackjack.BlackjackGameRules, hands int64, bankrole float32) GameResults {
	deck := blackjack.GenerateShoe(6).Shuffle()

	totalGames := 0
	playedHands := 0
	netWins := 0
	netLosses := 0
	blackjacks := 0

	for i := 0; i < handsPerGame; i++ {
		var handResults []blackjack.HandResult
		before := bankrole
		handResults, bankrole = blackjack.PlayHand(deck, rules, bankrole)
		if bankrole > before {
			netWins++
		} else if bankrole < before {
			netLosses++
		}
		if bankrole <= 0 {
			break
		}

		if handResults[0] == blackjack.HandResultBlackjack {
			blackjacks++
		}

		if deck.Remaining() < blackjack.DeckSize*1.5 {
			deck.Shuffle()
		}
		playedHands += len(handResults)
		totalGames++
	}

	return GameResults{
		Result:     bankrole,
		Hands:      totalGames,
		Blackjacks: float32(blackjacks) / float32(totalGames) * 100,
		WinRate:    float32(netWins) / float32(totalGames) * 100,
		LoseRate:   float32(netLosses) / float32(totalGames) * 100,
		PushRate:   float32(totalGames-netWins-netLosses) / float32(totalGames) * 100,
		EV:         (bankrole - startingBankrole) / float32(playedHands),
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	bjRules := blackjack.NewBlackjackGameRules(blackjack.InitGame(blackjack.H17Rules, blackjack.H17Splits))
	bjRules.SetDealerHitsSoft17(false)
	bjRules.SetDoubleAfterSplit(true) // not implemented
	bjRules.SetMaxPlayerSplits(2)
	bjRules.SetUseSimpleDeviations(true)
	bjRules.SetUseHighLowCounting(false) // <---- Enable/disable card counting
	bjRules.SetBidspread(blackjack.NewBidspread(
		map[int]blackjack.BidStrategy{
			0: {Hands: 1, Units: 1},
			1: {Hands: 1, Units: 2},
			2: {Hands: 1, Units: 4},
			3: {Hands: 1, Units: 8},
			4: {Hands: 1, Units: 12},
		}))

	resultsChannel := make(chan GameResults)
	overallResults := make([]GameResults, 0, iterations)
	sync := make(chan bool, 16)
	for i := 0; i < 16; i++ {
		sync <- true
	}

	for i := 0; i < iterations; i++ {
		go func(idx int) {
			<-sync
			resultsChannel <- PlayGame(bjRules, handsPerGame, startingBankrole)
			sync <- true
		}(i)
	}

	resultIdx := 0
	for res := range resultsChannel {
		overallResults = append(overallResults, res)
		log.Printf("%d: $%f | %+v", resultIdx, res.Result, res)
		resultIdx++
		if resultIdx >= iterations {
			break
		}
	}

	aggregatedResults := GameResults{}

	log.Println("====================================")
	for _, v := range overallResults {
		aggregatedResults.EV += v.EV
		aggregatedResults.Hands += v.Hands
		aggregatedResults.Blackjacks += v.Blackjacks
		aggregatedResults.WinRate += v.WinRate
		aggregatedResults.LoseRate += v.LoseRate
		aggregatedResults.PushRate += v.PushRate
		aggregatedResults.Result += v.Result
	}

	aggregatedResults.EV /= float32(iterations)
	aggregatedResults.WinRate /= float32(iterations)
	aggregatedResults.LoseRate /= float32(iterations)
	aggregatedResults.PushRate /= float32(iterations)
	aggregatedResults.Blackjacks /= float32(iterations)
	aggregatedResults.Result /= float32(iterations)

	log.Printf("Overall: %+v", aggregatedResults)
	// for k, v := range results {
	// 	log.Printf("\t%s:\t\t%f", ResultToString(k), float32(v)/float32(totalGames)*100)
	// }
}
