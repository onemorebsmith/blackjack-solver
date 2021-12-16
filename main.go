package main

import (
	"log"
	"math/rand"
	"time"
)

const startingBankrole = float32(1000) // in bet units
const handsPerGame = 10000
const iterations = 100

type GameResults struct {
	Hands    int
	WinRate  float32
	LoseRate float32
	PushRate float32
	EV       float32
	Result   float32
}

func PlayGame(rules *BlackjackGameRules, hands int64, bankrole float32) GameResults {
	deck := GenerateShoe(6)

	totalGames := 0
	playedHands := 0
	netWins := 0
	netLosses := 0

	// results := map[HandResult]int{
	// 	HandResultBlackjack:       0,
	// 	HandResultDealerBlackjack: 0,
	// 	HandResultLose:            0,
	// 	HandResultPush:            0,
	// 	HandResultWin:             0,
	// }
	for i := 0; i < handsPerGame; i++ {
		var handResults []HandResult
		before := bankrole
		handResults, bankrole = PlayHand(deck, rules, bankrole)
		if bankrole > before {
			netWins++
		} else if bankrole < before {
			netLosses++
		}
		if bankrole <= 0 {
			break
		}
		if deck.idx > DeckSize*5 {
			deck.Shuffle()
		}
		playedHands += len(handResults)
		totalGames++
	}

	return GameResults{
		Result:   bankrole,
		Hands:    totalGames,
		WinRate:  float32(netWins) / float32(totalGames) * 100,
		LoseRate: float32(netLosses) / float32(totalGames) * 100,
		PushRate: float32(totalGames-netWins-netLosses) / float32(totalGames) * 100,
		EV:       (bankrole - startingBankrole) / float32(playedHands),
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	bjRules := NewBlackjackGameRules(InitGame(h17Rules, h17Splits))
	bjRules.SetDealerHitsSoft17(true)
	bjRules.SetDoubleAfterSplit(true)
	bjRules.SetMaxPlayerSplits(2)
	bjRules.SetUseSimpleDeviations(true)
	bjRules.SetBidspread(NewBidspread(
		map[int]float32{
			0: 1,
			1: 1,
			2: 3,
			3: 5,
			4: 8,
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
		aggregatedResults.WinRate += v.WinRate
		aggregatedResults.LoseRate += v.LoseRate
		aggregatedResults.PushRate += v.PushRate
		aggregatedResults.Result += v.Result
	}

	aggregatedResults.EV /= float32(iterations)
	aggregatedResults.WinRate /= float32(iterations)
	aggregatedResults.LoseRate /= float32(iterations)
	aggregatedResults.PushRate /= float32(iterations)
	aggregatedResults.Result /= float32(iterations)

	log.Printf("Overall: %+v", aggregatedResults)
	// for k, v := range results {
	// 	log.Printf("\t%s:\t\t%f", ResultToString(k), float32(v)/float32(totalGames)*100)
	// }
}
