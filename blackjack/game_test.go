package blackjack

import (
	"testing"

	"github.com/onemorebsmith/blackjack-solver/blackjack/core"
)

func makeTestDeck() *core.Deck {
	return &core.Deck{
		Cards: []core.Card{
			{Value: 10},
			{Value: 3},
			{Value: 5},
			{Value: 11},
		},
	}
}

func PlayDealerHand(t *testing.T, dealerHand Hand, rules *BlackjackGameRules) Hand {
	if rules == nil {
		rules = MakeTestRules()
	}
	deck := makeTestDeck()
	return rules.PlayDealerHand(dealerHand, deck)
}

func PlaySingleTestHand(t *testing.T, playerHand Hand, dealerUpcard int) []Hand {
	rules := MakeTestRules()
	deck := makeTestDeck()
	return rules.PlayPlayerHand(playerHand, core.Card{Value: dealerUpcard}, deck, 1, 0)
}

func PlaySingleTestHandNonSplit(t *testing.T, playerHand Hand, dealerUpcard int) Hand {
	res := PlaySingleTestHand(t, playerHand, dealerUpcard)
	if len(res) != 1 {
		t.Fatalf("Should not have split %s vs %d", playerHand.toString(), dealerUpcard)
	}
	return res[0]
}

func Test_PlayNaturalBlackjackHand(t *testing.T) {
	res := PlaySingleTestHandNonSplit(t, MakeHand(11, 10), 2)
	if len(res.Cards) != 2 {
		t.Fatalf("Should not have hit on a blackjack")
	}
	handRes, result := CalculateHandResult(res, 18, 10)
	if handRes != HandResultBlackjack {
		t.Fatalf("result should have been blackjack")
	}
	if result != 15 {
		t.Fatalf("incorrect blackjack payout")
	}
}

func Test_PlayH12Hand(t *testing.T) {
	// hard 12 v dealer 8. Should hit (and bust)
	res := PlaySingleTestHandNonSplit(t, MakeHand(8, 4), 8)
	if len(res.Cards) != 3 {
		t.Fatalf("Should have hit on a 12")
	}
	if v, _ := res.HandValue(); v != 22 {
		t.Fatalf("Hand value should be 22, got %d", v)
	}
	handRes, result := CalculateHandResult(res, 18, 10)
	if handRes != HandResultLose {
		t.Fatalf("result should have been bust")
	}
	if result != -10 {
		t.Fatalf("incorrect losing payout")
	}
}

func Test_PlayH14Hand(t *testing.T) {
	// hard 14 v dealer 2. Should stand
	res := PlaySingleTestHandNonSplit(t, MakeHand(10, 4), 2)
	if len(res.Cards) != 2 {
		t.Fatalf("Should not hit on a 14 v 2")
	}
	if v, _ := res.HandValue(); v != 14 {
		t.Fatalf("Hand value should be 14, got %d", v)
	}
	handRes, result := CalculateHandResult(res, 22, 10)
	if handRes != HandResultWin {
		t.Fatalf("result should have been win")
	}
	if result != 10 {
		t.Fatalf("incorrect losing payout")
	}
}

func Test_PlayPushHand(t *testing.T) {
	// hard 20 v dealer 10. Should stand
	res := PlaySingleTestHandNonSplit(t, MakeHand(10, 10), 10)
	if len(res.Cards) != 2 {
		t.Fatalf("Should not hit on a 20 v 10")
	}
	if v, _ := res.HandValue(); v != 20 {
		t.Fatalf("Hand value should be 20, got %d", v)
	}
	handRes, result := CalculateHandResult(res, 20, 10)
	if handRes != HandResultPush {
		t.Fatalf("result should have been win")
	}
	if result != 0 {
		t.Fatalf("incorrect losing payout")
	}
}

func Test_Play88Hand(t *testing.T) {
	// 88 v dealer 6. Should split
	res := PlaySingleTestHand(t, MakeHand(8, 8), 6)
	if len(res) != 2 {
		t.Fatalf("Should have split on a double 8s")
	}
	if v, _ := res[0].HandValue(); v != 18 { // 8 10
		t.Fatalf("Hand 1 should be 18, got %d", v)
	}
	if v, _ := res[1].HandValue(); v != 16 { // 8 10
		t.Fatalf("Hand 2 should be 16, got %d", v)
	}
}

func Test_Play888Hand(t *testing.T) {
	deck := &core.Deck{
		Cards: []core.Card{
			{Value: 8},
			{Value: 8}, // hand 1 should split twice, then get 8/3 and double
			{Value: 3},
			{Value: 4},
			{Value: 4},
			{Value: 5}, // hand 2 should get 8/4/5 and stand
			{Value: 6},
			{Value: 7}, // hand 3 should get 8/6/7
			{Value: 11},
			{Value: 10}, // hand 4 should get 8A
		},
	}
	// 88 v dealer 6. Should split
	rules := MakeTestRules().SetMaxPlayerSplits(3)
	cards := MakeHand(8, 8)
	res := rules.PlayPlayerHand(cards, core.Card{Value: 7}, deck, 1, 0)
	if len(res) != 4 {
		t.Fatalf("Should have split on a double 8s")
	}
	if v, _ := res[0].HandValue(); v != 15 { // 8 10
		t.Fatalf("Hand 1 should be 15, got %d", v)
	}
	if res[0].Doubled == false { // 8 10
		t.Fatalf("Hand 1 have doubled on 11")
	}
	if v, _ := res[1].HandValue(); v != 17 { // 8 10
		t.Fatalf("Hand 2 should be 17, got %d", v)
	}
	if v, _ := res[2].HandValue(); v != 21 { // 8 10
		t.Fatalf("Hand 3 should be 21, got %d", v)
	}
	if v, _ := res[3].HandValue(); v != 19 { // 8 10
		t.Fatalf("Hand 4 should be 19, got %d", v)
	}
}

func Test_PlayAAAHand(t *testing.T) {
	// AA v dealer 6. Should split
	deck := &core.Deck{
		Cards: []core.Card{
			{Value: 11}, // P
			{Value: 10}, // D
			{Value: 11}, // P
			{Value: 6},  // D
			{Value: 10}, // hand 1 should be A10 -> 21
			{Value: 11}, // hand 2 should be AA -> 12
			{Value: 7},  // dealer busts w/ 23
		},
	}
	result, newBankrole := PlayHand(deck, MakeTestRules())
	_ = result
	if newBankrole != 2 {
		t.Fatalf("did not properly pay out after split aces, expected 2 got %f", newBankrole)
	}
}

func Test_DealerDoubleAces(t *testing.T) {
	rules := MakeTestRules()
	cards := MakeHand(11, 11)
	deck := makeTestDeck()
	res := rules.PlayDealerHand(cards, deck)
	if len(res.Cards) != 5 { // A A 10 3 5
		t.Fatalf("Dealer should hit on double aces")
	}
	if v, _ := res.HandValue(); v != 20 {
		t.Fatalf("Dealer hand should be 20, got %d", v)
	}
}

func Test_DealerH17(t *testing.T) {
	res := PlayDealerHand(t, MakeHand(6, 11), MakeTestRules().SetDealerHitsSoft17(false))
	if len(res.Cards) != 2 {
		t.Fatalf("H17 Rules: Dealer should stand on soft 17")
	}
	if v, _ := res.HandValue(); v != 17 {
		t.Fatalf("Dealer hand should be 17, got %d", v)
	}
}

func Test_DealerS17(t *testing.T) {
	res := PlayDealerHand(t, MakeHand(6, 11), MakeTestRules().SetDealerHitsSoft17(true))
	if len(res.Cards) != 3 {
		t.Fatalf("S17 Rules: Dealer should hit on soft 17")
	}
	if v, _ := res.HandValue(); v != 17 {
		t.Fatalf("Dealer hand should be 17, got %d", v)
	}
}

func Test_DealerHits(t *testing.T) {
	res := PlayDealerHand(t, MakeHand(2, 2), MakeTestRules())
	if len(res.Cards) != 4 {
		t.Fatalf("S17 Rules: Dealer should hit on soft 17")
	}
	if v, _ := res.HandValue(); v != 17 {
		t.Fatalf("Dealer hand should be 17, got %d", v)
	}
}
