package blackjack

import (
	"github.com/onemorebsmith/blackjack-solver/blackjack/core"
	"github.com/onemorebsmith/blackjack-solver/blackjack/strategies"
)

const blackjackPayout = float32(1.5)

type BlackjackGameRules struct {
	playerStrategy   *Ruleset
	DealerHitsSoft17 bool
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

func PlayHand(d *core.Deck, rules *BlackjackGameRules) ([]HandResult, float32) {
	bidStrategy := rules.TrackingStrategy.Bid(*d)
	perHandBid := bidStrategy.Units
	totalBid := bidStrategy.Units * float32(bidStrategy.Hands)
	//bankrole -= bidStrategy.Units * float32(bidStrategy.Hands)
	playerCards := Hand{}
	dealerCards := Hand{}

	playerCards.Cards = append(playerCards.Cards, d.Deal())
	dealerCards.Cards = append(dealerCards.Cards, d.Deal())
	playerCards.Cards = append(playerCards.Cards, d.Deal())
	dealerCards.Cards = append(dealerCards.Cards, d.Deal())
	dealerUpcard := dealerCards.Cards[1]

	insurance := false
	if dealerUpcard.Value == 11 { // insurance?
		// if rules.UseHighLowCounting {
		// 	tc := d.TrueCount()
		// 	if rules.UseSimpleDeviations && tc >= 3 {
		// 		insurance = true
		// 	}
		// }
	}

	// Check for dealer natural 21
	if dealerValue, _ := dealerCards.HandValue(); dealerValue == 21 {
		if playerValue, _ := playerCards.HandValue(); playerValue == 21 {
			// push if player has natural 21 also
			return []HandResult{HandResultBlackjackPush}, 0
		}
		if insurance {
			return []HandResult{HandResultInsuranceSave}, 0
		} else {
			return []HandResult{HandResultDealerBlackjack}, -totalBid
		}
	}

	playerHands := rules.PlayPlayerHand(playerCards, dealerUpcard, d, bidStrategy.Units, 0)
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

	dealerValue, _ := dealerCards.HandValue()
	results := []HandResult{}
	totalChange := float32(0)
	for _, h := range playerHands {
		handResult, change := CalculateHandResult(h, dealerValue, perHandBid)
		//log.Printf("%s ev: %f", handResult.ToString(), change)
		totalChange += change
		results = append(results, handResult)
	}
	//log.Println("--------------------")
	return results, totalChange
}

func CalculateHandResult(h Hand, dealerValue int, bid float32) (HandResult, float32) {
	playerValue, _ := h.HandValue()
	playerBlackjack := playerValue == 21
	playerNaturalBlackjack := playerBlackjack && len(h.Cards) == 2 && !h.SplitHand

	if playerNaturalBlackjack {
		return HandResultBlackjack, (bid * blackjackPayout)
	}
	if h.Doubled {
		bid *= 2
	}
	if playerValue == dealerValue {
		// push
		return HandResultPush, 0
	}
	if playerValue > 21 {
		// player busted, instant loss
		return HandResultLose, -bid
	}
	if (dealerValue > 21) || playerValue > dealerValue {
		// player wins
		return HandResultWin, bid
	}

	// dealer wins
	return HandResultLose, -bid
}

func (rs *BlackjackGameRules) PlayPlayerHand(playerHand Hand, dealerUpcard core.Card,
	deck *core.Deck, bid float32, splitCounter int) []Hand {
	//log.Printf("%s vs %s", playerHand.toString(), dealerUpcard.ToString())
	finished := false
	if bid > 5 {
		finished = false
	}
	for {
		decision := rs.MakePlayerDecision(playerHand, dealerUpcard, splitCounter)
		switch decision {
		case PlayerDecisionNatural21:
			finished = true
		case PlayerDecisionDouble:
			if playerHand.SplitHand {
				finished = true
			}
			playerHand.Cards = append(playerHand.Cards, deck.Deal())
			playerHand.Doubled = true
			finished = true
		case PlayerDecisionSplitAces:
			// TODO(bs): handle RSA here
			handA := Hand{Cards: []core.Card{playerHand.Cards[0], deck.Deal()}, SplitHand: true}
			handB := Hand{Cards: []core.Card{playerHand.Cards[1], deck.Deal()}, SplitHand: true}
			return []Hand{handA, handB} // can only take one card after aces
		case PlayerDecisionSplit:
			splitCounter++
			handA := rs.PlayPlayerHand(Hand{Cards: []core.Card{playerHand.Cards[0], deck.Deal()}, SplitHand: true}, dealerUpcard, deck, bid, splitCounter)
			handB := rs.PlayPlayerHand(Hand{Cards: []core.Card{playerHand.Cards[1], deck.Deal()}, SplitHand: true}, dealerUpcard, deck, bid, splitCounter)
			return append(handA, handB...)
		case PlayerDecisionHit:
			playerHand.Cards = append(playerHand.Cards, deck.Deal())
		case PlayerDecisionStand:
			//log.Printf("  stand %s", playerHand.toString())
			finished = true
		}

		v, _ := playerHand.HandValue()
		if v >= 21 || finished {
			break
		}
	}

	return []Hand{playerHand}
}

func (rs *BlackjackGameRules) PlayDealerHand(dealerHand Hand, deck *core.Deck) Hand {
	for {
		decision := rs.MakeDealerDecision(dealerHand)
		if decision == PlayerDecisionHit {
			newCard := deck.Deal()
			dealerHand.Cards = append(dealerHand.Cards, newCard)
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
	tcMap := map[int]int{}
	deck.PreviewCard = func(c core.Card) {
		rules.TrackingStrategy.Update(c)
	}
	totalHands := 0
	for {
		totalHands++
		handResults, change := PlayHand(deck, rules)
		if change > 0 {
			netWins++
		} else if bankrole < before {
			netLosses++
		}
		bankrole += change
		if bankrole <= 0 {
			break
		}

		for _, hr := range handResults {
			if hr == HandResultBlackjack {
				blackjacks++
			}
		}

		//tcMap[deck.Count]++
		if deck.Remaining() < int(core.DeckSize*rules.Penetration) {
			//deck.Shuffle()
			break
		}
	}
	aggregatedTC := 0
	for count, frequency := range tcMap {
		aggregatedTC += count * frequency
		totalHands += frequency
	}
	averageTC := float32(aggregatedTC) / float32(totalHands)

	return GameResults{
		Result:     bankrole,
		Hands:      totalHands,
		Blackjacks: blackjacks,
		Wins:       netWins,
		Losses:     netLosses,
		Pushes:     totalHands - netWins - netLosses,
		EV:         bankrole - before,
		AvgTc:      averageTC,
		TCMap:      tcMap,
	}
}
