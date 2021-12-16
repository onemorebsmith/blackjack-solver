package main

import "log"

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
		bidSpread:           NewBidspread(map[int]float32{0: 1}),
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
	bid := float32(1)
	if rules.UseHighLowCounting {
		bid = rules.bidSpread.Bid(d)
	}
	if bankrole < bid {
		return nil, bankrole
	}

	playerCards := Hand{}
	dealerCards := Hand{}

	playerCards.Cards = append(playerCards.Cards, d.Deal())
	dealerCards.Cards = append(dealerCards.Cards, d.Deal())
	playerCards.Cards = append(playerCards.Cards, d.Deal())
	dealerCards.Cards = append(dealerCards.Cards, d.Deal())

	dealerUpcard := dealerCards.Cards[0].Value
	insurance := false
	if dealerUpcard == 11 { // insurance?
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

	playerHands := rules.PlayPlayerHand(playerCards, dealerCards.Cards[0], d, 0)
	dealerCards = rules.PlayDealerHand(dealerCards, d)

	dealerValue, _ := dealerCards.HandValue()
	results := []HandResult{}
	var handResult HandResult
	for _, h := range playerHands {
		handResult, bankrole = CalculateHandResult(h, dealerValue, bankrole, bid)
		results = append(results, handResult)
	}
	return results, bankrole
}

func CalculateHandResult(h Hand, dealerValue int, bankrole float32, bid float32) (HandResult, float32) {
	if h.Doubled {
		bid *= 2
	}
	playerValue, _ := h.HandValue()
	if playerValue > 21 {
		// player busted
		return HandResultLose, bankrole - bid
	} else if playerValue == 21 && len(h.Cards) == 2 && !h.SplitHand {
		// natural blackjack, can't happen on splits
		return HandResultBlackjack, bankrole + (bid * blackjackPayout)
	} else if dealerValue > 21 && playerValue <= 21 {
		// dealer busted
		return HandResultWin, bankrole + bid
	} else if dealerValue == playerValue {
		// push
		return HandResultPush, bankrole
	} else if dealerValue > playerValue {
		return HandResultLose, bankrole - bid
	} else if playerValue > dealerValue { // player Win
		return HandResultWin, bankrole + bid
	} else {
		log.Println("Error? unknown hand result")
		return HandResultPush, bankrole
	}
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
		dealerValue, soft := dealerHand.HandValue()
		if rs.DealerHitsSoft17 && soft && dealerValue == 17 { // H17 vs S17 rule
			dealerHand.Cards = append(dealerHand.Cards, deck.Deal())
		} else if dealerValue < 17 {
			dealerHand.Cards = append(dealerHand.Cards, deck.Deal())
		} else {
			break
		}
	}
	return dealerHand
}
