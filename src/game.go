package blackjack

import (
	"math"

	"github.com/onemorebsmith/blackjack-solver/src/blackjack/core"
	"github.com/onemorebsmith/blackjack-solver/src/blackjack/strategies"
)

const blackjackPayout = float32(1.5)

type BlackjackGameRules struct {
	playerStrategy   *Ruleset
	DealerHitsSoft17 bool
	ReSplitAces      bool
	MaxPlayerSplits  int
	DoubleAfterSplit bool
	Penetration      float32
	TrackingStrategy strategies.TrackingStrategy

	UseSimpleDeviations bool // use insurance after TC 3+ & no hit 12
}

func NewBlackjackGameRules(rules *Ruleset) *BlackjackGameRules {
	return &BlackjackGameRules{
		playerStrategy:      rules,
		DealerHitsSoft17:    true,
		MaxPlayerSplits:     4,
		DoubleAfterSplit:    true,
		UseSimpleDeviations: false,
		ReSplitAces:         false,
		TrackingStrategy:    strategies.InitFlatbetStrategy(),
	}
}

func (bj *BlackjackGameRules) SetPenetration(pen float32) *BlackjackGameRules {
	bj.Penetration = pen
	return bj
}

func (bj *BlackjackGameRules) SetDealerHitsSoft17(v bool) *BlackjackGameRules {
	bj.DealerHitsSoft17 = v
	return bj
}

func (bj *BlackjackGameRules) SetResplitAces(v bool) *BlackjackGameRules {
	bj.ReSplitAces = v
	return bj
}

func (bj *BlackjackGameRules) SetMaxPlayerSplits(v int) *BlackjackGameRules {
	bj.MaxPlayerSplits = v
	return bj
}

func (bj *BlackjackGameRules) SetDoubleAfterSplit(v bool) *BlackjackGameRules {
	bj.DoubleAfterSplit = v
	return bj
}

func (bj *BlackjackGameRules) SetUseSimpleDeviations(v bool) *BlackjackGameRules {
	bj.UseSimpleDeviations = v
	return bj
}

func PlayGame(rules BlackjackGameRules, decks int, shoes int, bankrole float32, handsPerHour float32) GameResults {
	deck := core.GenerateShoe(decks).Shuffle()
	// create a new instance of the tracking strategy as to not share state
	// with the other threads
	rules.TrackingStrategy = rules.TrackingStrategy.Instance()

	deck.PreviewCard = func(c core.Card) {
		rules.TrackingStrategy.Update(c)
	}
	totalGames := 0

	handAVs := make([]float32, 0, shoes*50) // shoes average ~45 hands heads up
	aggregatedResults := GameResults{}
	for i := 0; i < shoes; i++ {
		result := PlayShoe(deck, &rules, bankrole)
		if hl, ok := rules.TrackingStrategy.(*strategies.HighLowCountStrategy); ok {
			result.BidsByTC = hl.BidsByTC
			result.AvgTC = hl.AggregatedTC / float32(hl.Updates)
			result.HighTC = hl.HighTC
			result.LowTC = hl.LowTC
		}
		rules.TrackingStrategy.Shuffle()
		deck.Shuffle()
		totalGames++
		aggregatedResults = AggregateResults(aggregatedResults, result)
		handAVs = append(handAVs, result.HandAVs...)
	}
	// calculate the population standard dev
	evAgg := float32(0)
	handsGroupedHourly := []float32{}
	hourlyHandCounter := float32(0)
	hourlyAgg := float32(0)
	hourlyOverallTotal := float32(0)
	for _, av := range handAVs {
		evAgg += av
		hourlyAgg += av
		hourlyHandCounter++
		// this is truly horrendous and should be cleaned up, but in order to calculate hourly standard
		// dev you need to calculate standard dev across the aggregated hourly AVs instead of
		// the individual hand AVs
		if hourlyHandCounter > handsPerHour {
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

func PlayHand(d *core.Deck, rules *BlackjackGameRules) []core.HandResult {
	bidStrategy := rules.TrackingStrategy.Bid(*d)
	perHandBid := bidStrategy.Units
	playerCards := core.Hand{}
	dealerCards := core.Hand{}
	// TODO (bs): support multiple hands
	playerCards.Cards = append(playerCards.Cards, d.Deal())
	dealerCards.Cards = append(dealerCards.Cards, d.Deal())
	playerCards.Cards = append(playerCards.Cards, d.Deal())
	dealerCards.Cards = append(dealerCards.Cards, d.Deal())
	dealerUpcard := dealerCards.Cards[1]

	//insurance := false
	if dealerUpcard.Value == 11 { // insurance?
		// if rules.UseHighLowCounting {
		// 	tc := d.TrueCount()
		// 	if rules.UseSimpleDeviations && tc >= 3 {
		// 		insurance = true
		// 	}
		// }
	}

	playerHands := []core.Hand{playerCards}
	// Play the hand if the dealer does not have 21
	if dealerValue, _ := dealerCards.HandValue(); dealerValue != 21 {
		splitCounter := 0
		playerHands = rules.PlayPlayerHand(playerCards, dealerUpcard, d, bidStrategy.Units, &splitCounter)
		allBusted := true
		for _, v := range playerHands {
			if handVal, _ := v.HandValue(); handVal <= 21 {
				allBusted = false
				break
			}
		}
		if !allBusted {
			dealerCards = rules.PlayDealerHand(dealerCards, d)
		}
	}

	results := make([]core.HandResult, 0, len(playerHands))
	totalChange := float32(0)
	for _, h := range playerHands {
		handResult := CalculateHandResult(h, dealerCards, perHandBid)
		totalChange += handResult.AV
		results = append(results, handResult)
	}
	return results
}

func CalculateHandResult(playerHand core.Hand, dealerHand core.Hand, bid float32) core.HandResult {
	if playerHand.Doubled {
		bid *= 2
	}

	playerValue, _ := playerHand.HandValue()
	dealerValue, _ := dealerHand.HandValue()
	playerNaturalBlackjack := playerValue == 21 && playerHand.IsNatural()
	dealerNaturalBlackjack := dealerValue == 21 && dealerHand.IsNatural()

	// deal with blackjacks
	if dealerNaturalBlackjack && playerNaturalBlackjack {
		return core.MakeHandResult(core.HandResultBlackjackPush, 0)
	} else if dealerNaturalBlackjack {
		return core.MakeHandResult(core.HandResultDealerBlackjack, -bid)
	} else if playerNaturalBlackjack {
		return core.MakeHandResult(core.HandResultBlackjack, bid*blackjackPayout)
	}

	if playerValue == dealerValue { // push
		return core.HandResult{Result: core.HandResultPush, AV: 0}
	} else if playerValue > 21 { // player bust
		// player busted, instant loss
		return core.HandResult{Result: core.HandResultLose, AV: -bid}
	} else if dealerValue > 21 { // dealer bust
		return core.HandResult{Result: core.HandResultWin, AV: bid}
	} else if playerValue > dealerValue { // player wins
		return core.HandResult{Result: core.HandResultWin, AV: bid}
	} else { // dealer wins
		return core.HandResult{Result: core.HandResultLose, AV: -bid}
	}
}

func (rs *BlackjackGameRules) PlayPlayerHand(playerHand core.Hand, dealerUpcard core.Card,
	deck *core.Deck, bid float32, splitCounter *int) []core.Hand {
	finished := false
	for {
		decision := rs.MakePlayerDecision(playerHand, dealerUpcard, *splitCounter)
		switch decision {
		case PlayerDecisionNatural21:
			finished = true
		case PlayerDecisionDouble:
			playerHand.Cards = append(playerHand.Cards, deck.Deal())
			playerHand.Doubled = true
			finished = true
		case PlayerDecisionSplitAces:
			hands := make([]core.Hand, 0, 4)
			hands = append(hands, rs.PlayPlayerHand(core.Hand{Cards: []core.Card{playerHand.Cards[0], deck.Deal()}, SplitAcesHand: true},
				dealerUpcard, deck, bid, splitCounter)...)
			hands = append(hands, rs.PlayPlayerHand(core.Hand{Cards: []core.Card{playerHand.Cards[1], deck.Deal()}, SplitAcesHand: true},
				dealerUpcard, deck, bid, splitCounter)...)
			return hands // can only take one card after aces
		case PlayerDecisionSplit:
			*splitCounter++
			hands := make([]core.Hand, 0, 4)
			hands = append(hands, rs.PlayPlayerHand(core.Hand{Cards: []core.Card{playerHand.Cards[0], deck.Deal()}, SplitHand: true},
				dealerUpcard, deck, bid, splitCounter)...)
			hands = append(hands, rs.PlayPlayerHand(core.Hand{Cards: []core.Card{playerHand.Cards[1], deck.Deal()}, SplitHand: true},
				dealerUpcard, deck, bid, splitCounter)...)
			return hands
		case PlayerDecisionHit:
			playerHand.Cards = append(playerHand.Cards, deck.Deal())
		case PlayerDecisionStand:
			finished = true
		}
		v, _ := playerHand.HandValue()
		if v >= 21 || finished {
			break
		}
	}

	return []core.Hand{playerHand}
}

func (rs *BlackjackGameRules) PlayDealerHand(dealerHand core.Hand, deck *core.Deck) core.Hand {
	for {
		decision := rs.MakeDealerDecision(dealerHand)
		if decision == PlayerDecisionHit {
			dealerHand.Cards = append(dealerHand.Cards, deck.Deal())
		} else {
			break
		}
	}
	return dealerHand
}

func PlayShoe(deck *core.Deck, rules *BlackjackGameRules, bankrole float32) GameResults {
	before := bankrole
	netWins := 0
	netLosses := 0
	blackjacks := 0
	totalHands := 0
	handAVs := make([]float32, 0, 50)
	for {
		totalHands++
		handResults := PlayHand(deck, rules)
		handAV := float32(0)
		for _, r := range handResults {
			bankrole += r.AV
			handAV += r.AV
			if r.AV > 0 {
				netWins++
			} else if bankrole < before {
				netLosses++
			}
			if r.Result == core.HandResultBlackjack {
				blackjacks++
			}
		}
		handAVs = append(handAVs, handAV)
		if bankrole <= 0 {
			break
		}
		if deck.Remaining() < int(core.DeckSize*rules.Penetration) {
			break
		}
	}
	return GameResults{
		Result:     bankrole,
		Hands:      totalHands,
		Blackjacks: blackjacks,
		Wins:       netWins,
		Losses:     netLosses,
		Pushes:     totalHands - netWins - netLosses,
		EV:         bankrole - before,
		HandAVs:    handAVs,
	}
}
