package main

import (
	"testing"
)

func makeTestDeck() *Deck {
	return &Deck{
		idx: 0,
		Cards: []Card{
			{Value: 10},
			{Value: 3},
			{Value: 5},
			{Value: 11},
		},
	}
}

func Test_PlayNaturalBlackjackHand(t *testing.T) {
	rules := MakeTestRules()
	deck := makeTestDeck()

	cards := MakeHand(11, 10)
	res := rules.PlayPlayerHand(cards, Card{Value: 2}, deck, 0)
	if len(res) != 1 {
		t.Fatalf("Should not have doubled on a blackjack")
	}
	if len(res[0].Cards) != 2 {
		t.Fatalf("Should not have hit on a blackjack")
	}
	handRes, bankrole := CalculateHandResult(res[0], 18, 100, 10)
	if handRes != HandResultBlackjack {
		t.Fatalf("result should have been blackjack")
	}
	if bankrole != 115 {
		t.Fatalf("incorrect blackjack payout")
	}
}

func Test_PlayH12Hand(t *testing.T) {
	// hard 12 v dealer 8. Should hit (and bust)
	rules := MakeTestRules()
	deck := makeTestDeck()

	cards := MakeHand(8, 4)
	res := rules.PlayPlayerHand(cards, Card{Value: 8}, deck, 0)
	if len(res) != 1 {
		t.Fatalf("Should not have split")
	}
	if len(res[0].Cards) != 3 {
		t.Fatalf("Should have hit on a 12")
	}
	if v, _ := res[0].HandValue(); v != 22 {
		t.Fatalf("Hand value should be 22, got %d", v)
	}
	handRes, bankrole := CalculateHandResult(res[0], 18, 100, 10)
	if handRes != HandResultLose {
		t.Fatalf("result should have been bust")
	}
	if bankrole != 90 {
		t.Fatalf("incorrect losing payout")
	}
}

func Test_PlayH14Hand(t *testing.T) {
	// hard 14 v dealer 2. Should stand
	rules := MakeTestRules()
	deck := makeTestDeck()

	cards := MakeHand(10, 4)
	res := rules.PlayPlayerHand(cards, Card{Value: 2}, deck, 0)
	if len(res) != 1 {
		t.Fatalf("Should not have split")
	}
	if len(res[0].Cards) != 2 {
		t.Fatalf("Should not hit on a 14 v 2")
	}
	if v, _ := res[0].HandValue(); v != 14 {
		t.Fatalf("Hand value should be 14, got %d", v)
	}
	handRes, bankrole := CalculateHandResult(res[0], 22, 100, 10)
	if handRes != HandResultWin {
		t.Fatalf("result should have been win")
	}
	if bankrole != 110 {
		t.Fatalf("incorrect losing payout")
	}
}

func Test_PlayPushHand(t *testing.T) {
	// hard 14 v dealer 2. Should stand
	rules := MakeTestRules()
	deck := makeTestDeck()

	cards := MakeHand(10, 10)
	res := rules.PlayPlayerHand(cards, Card{Value: 10}, deck, 0)
	if len(res) != 1 {
		t.Fatalf("Should not have split")
	}
	if len(res[0].Cards) != 2 {
		t.Fatalf("Should not hit on a 20 v 10")
	}
	if v, _ := res[0].HandValue(); v != 20 {
		t.Fatalf("Hand value should be 20, got %d", v)
	}
	handRes, bankrole := CalculateHandResult(res[0], 20, 100, 10)
	if handRes != HandResultPush {
		t.Fatalf("result should have been win")
	}
	if bankrole != 100 {
		t.Fatalf("incorrect losing payout")
	}
}

func Test_Play88Hand(t *testing.T) {
	// AA v dealer 6. Should split
	rules := MakeTestRules()
	deck := makeTestDeck()

	cards := MakeHand(8, 8)
	res := rules.PlayPlayerHand(cards, Card{Value: 6}, deck, 0)
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

func Test_PlayAAAHand(t *testing.T) {
	// AA v dealer 6. Should split
	deck := &Deck{
		idx: 0,
		Cards: []Card{
			{Value: 10}, // hand 1 should be A10
			{Value: 11}, // hand 2 should be AA and split again
			{Value: 7},  // hand 2 should end up being A 7 4 (soft 18 -> doubled)
			{Value: 4},
			{Value: 5},
			{Value: 2}, // hand 3 should be A 5 2
		},
	}

	rules := MakeTestRules()
	cards := MakeHand(11, 11)
	res := rules.PlayPlayerHand(cards, Card{Value: 6}, deck, 0)
	if len(res) != 3 {
		t.Fatalf("Should have split on a 14 v 2")
	}
	if v, _ := res[0].HandValue(); v != 21 {
		t.Fatalf("Hand 1 should be 14, got %d", v)
	}
	if v, _ := res[1].HandValue(); v != 12 {
		t.Fatalf("Hand 2 should be 12, got %d", v)
	}
	if v, _ := res[2].HandValue(); v != 18 {
		t.Fatalf("Hand 2 should be 18, got %d", v)
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
