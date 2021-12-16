package blackjack

const blackjackPayout = float32(1.5)

type BlackjackGameRules struct {
	playerStrategy   *Ruleset
	bidSpread        *Bidspread
	DealerHitsSoft17 bool
	MaxPlayerSplits  int
	DoubleAfterSplit bool

	UseHighLowCounting  bool // true to enable high/low counting and bet spreading
	UseSimpleDeviations bool // use insurance after TC 3+ & no hit 12
}

func NewBlackjackGameRules(rules *Ruleset) *BlackjackGameRules {
	return &BlackjackGameRules{
		playerStrategy:      rules,
		DealerHitsSoft17:    false,
		MaxPlayerSplits:     2,
		DoubleAfterSplit:    true,
		UseSimpleDeviations: false,
		UseHighLowCounting:  false,
		bidSpread:           NewBidspread(map[int]BidStrategy{0: {Units: 1, Hands: 1}}),
	}
}

func (bj *BlackjackGameRules) SetBidspread(bs *Bidspread) *BlackjackGameRules {
	bj.bidSpread = bs
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

func (bj *BlackjackGameRules) SetUseHighLowCounting(v bool) *BlackjackGameRules {
	bj.UseHighLowCounting = v
	return bj
}

func PlayHand(d *Deck, rules *BlackjackGameRules, bankrole float32) ([]HandResult, float32) {
	bidStrategy := BidStrategy{Hands: 1, Units: 1}
	if rules.UseHighLowCounting {
		bidStrategy = rules.bidSpread.Bid(d)
	}

	bid := bidStrategy.Units
	if bankrole < bid*float32(bidStrategy.Hands) {
		return []HandResult{HandResultPush}, bankrole
	}

	playerCards := Hand{}
	dealerCards := Hand{}

	playerCards.Cards = append(playerCards.Cards, d.Deal())
	dealerCards.Cards = append(dealerCards.Cards, d.Deal())
	playerCards.Cards = append(playerCards.Cards, d.Deal())
	dealerCards.Cards = append(dealerCards.Cards, d.Deal())

	dealerUpcard := dealerCards.Cards[1]
	insurance := false
	if dealerUpcard.Value == 11 { // insurance?
		if rules.UseHighLowCounting {
			tc := d.TrueCount()
			if rules.UseSimpleDeviations && tc >= 3 {
				insurance = true
				bankrole -= (bid / 2)
			}
		}
	}

	// Check for dealer natural 21
	if dealerValue, _ := dealerCards.HandValue(); dealerValue == 21 {
		if playerValue, _ := playerCards.HandValue(); playerValue == 21 {
			// push if player has natural 21 also
			return []HandResult{HandResultPush}, bankrole
		}
		if insurance {
			return []HandResult{HandResultInsuranceSave}, bankrole + (bid / 2)
		} else {
			return []HandResult{HandResultDealerBlackjack}, bankrole - bid
		}
	}

	playerHands := rules.PlayPlayerHand(playerCards, dealerUpcard, d, 0)
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
	for _, h := range playerHands {
		handResult, change := CalculateHandResult(h, dealerValue, bid)
		bankrole += change
		results = append(results, handResult)
	}

	return results, bankrole
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

func (rs *BlackjackGameRules) PlayPlayerHand(playerHand Hand, dealerUpcard Card, deck *Deck, splitCounter int) []Hand {
	finished := false
	for {
		decision := rs.MakePlayerDecision(playerHand, dealerUpcard, splitCounter)
		switch decision {
		case PlayerDecisionDouble:
			playerHand.Cards = append(playerHand.Cards, deck.Deal())
			playerHand.Doubled = true
			finished = true
		case PlayerDecisionSplitAces:
			handA := Hand{Cards: []Card{playerHand.Cards[0], deck.Deal()}, SplitHand: true}
			handB := Hand{Cards: []Card{playerHand.Cards[1], deck.Deal()}, SplitHand: true}
			return []Hand{handA, handB} // can only take one card after aces
		case PlayerDecisionSplit:
			splitCounter++
			handA := rs.PlayPlayerHand(Hand{Cards: []Card{playerHand.Cards[0], deck.Deal()}, SplitHand: true}, dealerUpcard, deck, splitCounter)
			if len(handA) > 1 {
				splitCounter++
			}
			handB := rs.PlayPlayerHand(Hand{Cards: []Card{playerHand.Cards[1], deck.Deal()}, SplitHand: true}, dealerUpcard, deck, splitCounter)
			return append(handA, handB...)
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

	return []Hand{playerHand}
}

func (rs *BlackjackGameRules) PlayDealerHand(dealerHand Hand, deck *Deck) Hand {
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
