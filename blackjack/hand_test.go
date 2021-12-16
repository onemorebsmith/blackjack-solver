package blackjack

import "testing"

func TestHandCounts(t *testing.T) {
	tests := []struct {
		Hand     Hand
		Expected int
		Soft     bool
	}{
		// Soft Ace hands
		{Hand: MakeHand(11, 11), Soft: true},
		{Hand: MakeHand(11, 2), Expected: 13, Soft: true},
		{Hand: MakeHand(11, 3), Expected: 14, Soft: true},
		{Hand: MakeHand(11, 4), Expected: 15, Soft: true},
		{Hand: MakeHand(11, 5), Expected: 16, Soft: true},
		{Hand: MakeHand(11, 6), Expected: 17, Soft: true},
		{Hand: MakeHand(11, 7), Expected: 18, Soft: true},
		{Hand: MakeHand(11, 8), Expected: 19, Soft: true},
		{Hand: MakeHand(11, 9), Expected: 20, Soft: true},
		{Hand: MakeHand(11, 10), Expected: 21, Soft: true},

		{Hand: MakeHand(11, 11, 11, 11, 11), Expected: 15, Soft: true},
		{Hand: MakeHand(11, 8, 2), Expected: 21, Soft: true},

		// Hard Ace hands
		{Hand: MakeHand(11, 10, 10), Expected: 21, Soft: false}, // Ace is no longer soft, can only be a 1
		{Hand: MakeHand(11, 10, 9), Expected: 20, Soft: false},

		{Hand: MakeHand(10, 10), Expected: 20, Soft: false},
		{Hand: MakeHand(10, 9), Expected: 19, Soft: false},
		{Hand: MakeHand(10, 6), Expected: 16, Soft: false},
		{Hand: MakeHand(2, 3, 4, 5, 7), Expected: 21, Soft: false},
	}

	for _, tc := range tests {
		val, soft := tc.Hand.HandValue()
		if soft != tc.Soft {
			t.Fatalf("Hand soft flag is incorrect, expected: %t, got: %t", tc.Soft, soft)
		}
		if val != tc.Expected {
			t.Fatalf("Hand value is incorrect, expected: %d, got: %d", tc.Expected, val)
		}
	}
}

func TestH17SoftHandDecisions(t *testing.T) {
	tests := []struct {
		Hand              Hand
		Expected          int
		ExpectedDecisions map[int]PlayerDecision
	}{
		// Soft Ace hands
		{Hand: MakeHand(11, 11), ExpectedDecisions: map[int]PlayerDecision{
			-1: PlayerDecisionSplitAces,
		}},

		// A2
		{Hand: MakeHand(11, 2), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionHit,
			3:  PlayerDecisionHit,
			4:  PlayerDecisionHit,
			5:  PlayerDecisionDouble,
			6:  PlayerDecisionDouble,
			7:  PlayerDecisionHit,
			8:  PlayerDecisionHit,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},

		// A3
		{Hand: MakeHand(11, 3), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionHit,
			3:  PlayerDecisionHit,
			4:  PlayerDecisionHit,
			5:  PlayerDecisionDouble,
			6:  PlayerDecisionDouble,
			7:  PlayerDecisionHit,
			8:  PlayerDecisionHit,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// A4
		{Hand: MakeHand(11, 4), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionHit,
			3:  PlayerDecisionHit,
			4:  PlayerDecisionDouble,
			5:  PlayerDecisionDouble,
			6:  PlayerDecisionDouble,
			7:  PlayerDecisionHit,
			8:  PlayerDecisionHit,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// A5
		{Hand: MakeHand(11, 5), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionHit,
			3:  PlayerDecisionHit,
			4:  PlayerDecisionDouble,
			5:  PlayerDecisionDouble,
			6:  PlayerDecisionDouble,
			7:  PlayerDecisionHit,
			8:  PlayerDecisionHit,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// A6
		{Hand: MakeHand(11, 6), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionHit,
			3:  PlayerDecisionDouble,
			4:  PlayerDecisionDouble,
			5:  PlayerDecisionDouble,
			6:  PlayerDecisionDouble,
			7:  PlayerDecisionHit,
			8:  PlayerDecisionHit,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// A7
		{Hand: MakeHand(11, 7), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionDouble,
			3:  PlayerDecisionDouble,
			4:  PlayerDecisionDouble,
			5:  PlayerDecisionDouble,
			6:  PlayerDecisionDouble,
			7:  PlayerDecisionStand,
			8:  PlayerDecisionStand,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// A8
		{Hand: MakeHand(11, 8), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionStand,
			3:  PlayerDecisionStand,
			4:  PlayerDecisionStand,
			5:  PlayerDecisionStand,
			6:  PlayerDecisionDouble,
			7:  PlayerDecisionStand,
			8:  PlayerDecisionStand,
			9:  PlayerDecisionStand,
			10: PlayerDecisionStand,
			11: PlayerDecisionStand,
		}},
		// A9
		{Hand: MakeHand(11, 9), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionStand,
			3:  PlayerDecisionStand,
			4:  PlayerDecisionStand,
			5:  PlayerDecisionStand,
			6:  PlayerDecisionStand,
			7:  PlayerDecisionStand,
			8:  PlayerDecisionStand,
			9:  PlayerDecisionStand,
			10: PlayerDecisionStand,
			11: PlayerDecisionStand,
		}},
	}

	rules := MakeTestRules()

	for _, tc := range tests {
		defAction, hasDefault := tc.ExpectedDecisions[-1]
		for dealerCard := 2; dealerCard < 11; dealerCard++ {
			decision := rules.MakePlayerDecision(tc.Hand, Card{Value: dealerCard}, 0)
			if hasDefault {
				if defAction != decision {
					t.Fatalf("Unexpected default player decision for `%s` vs `%d`", tc.Hand.toString(), dealerCard)
				}
			} else {
				expectedAction := tc.ExpectedDecisions[dealerCard]
				if expectedAction != decision {
					t.Fatalf("Unexpected player decision for `%s` vs `%d`, got %s, expected %s",
						tc.Hand.toString(), dealerCard, decision.ToString(), expectedAction.ToString())
				}
			}
		}
	}
}

func TestH17HardHandDecisions(t *testing.T) {
	tests := []struct {
		Hand              Hand
		ExpectedDecisions map[int]PlayerDecision
	}{
		// Soft Ace hands
		{Hand: MakeHand(10, 7), ExpectedDecisions: map[int]PlayerDecision{
			-1: PlayerDecisionStand,
		}},

		// 16
		{Hand: MakeHand(10, 6), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionStand,
			3:  PlayerDecisionStand,
			4:  PlayerDecisionStand,
			5:  PlayerDecisionStand,
			6:  PlayerDecisionStand,
			7:  PlayerDecisionHit,
			8:  PlayerDecisionHit,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// 15
		{Hand: MakeHand(10, 5), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionStand,
			3:  PlayerDecisionStand,
			4:  PlayerDecisionStand,
			5:  PlayerDecisionStand,
			6:  PlayerDecisionStand,
			7:  PlayerDecisionHit,
			8:  PlayerDecisionHit,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// 14
		{Hand: MakeHand(10, 4), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionStand,
			3:  PlayerDecisionStand,
			4:  PlayerDecisionStand,
			5:  PlayerDecisionStand,
			6:  PlayerDecisionStand,
			7:  PlayerDecisionHit,
			8:  PlayerDecisionHit,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// 13
		{Hand: MakeHand(10, 3), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionStand,
			3:  PlayerDecisionStand,
			4:  PlayerDecisionStand,
			5:  PlayerDecisionStand,
			6:  PlayerDecisionStand,
			7:  PlayerDecisionHit,
			8:  PlayerDecisionHit,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// 12
		{Hand: MakeHand(10, 2), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionHit,
			3:  PlayerDecisionHit,
			4:  PlayerDecisionStand,
			5:  PlayerDecisionStand,
			6:  PlayerDecisionStand,
			7:  PlayerDecisionHit,
			8:  PlayerDecisionHit,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// 11
		{Hand: MakeHand(9, 2), ExpectedDecisions: map[int]PlayerDecision{
			-1: PlayerDecisionDouble,
		}},
		// 10
		{Hand: MakeHand(5, 5), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionDouble,
			3:  PlayerDecisionDouble,
			4:  PlayerDecisionDouble,
			5:  PlayerDecisionDouble,
			6:  PlayerDecisionDouble,
			7:  PlayerDecisionDouble,
			8:  PlayerDecisionDouble,
			9:  PlayerDecisionDouble,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// 9
		{Hand: MakeHand(5, 4), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionHit,
			3:  PlayerDecisionDouble,
			4:  PlayerDecisionDouble,
			5:  PlayerDecisionDouble,
			6:  PlayerDecisionDouble,
			7:  PlayerDecisionHit,
			8:  PlayerDecisionHit,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// 8
		{Hand: MakeHand(3, 5), ExpectedDecisions: map[int]PlayerDecision{
			-1: PlayerDecisionHit,
		}},
	}

	rules := MakeTestRules()

	for _, tc := range tests {
		defAction, hasDefault := tc.ExpectedDecisions[-1]
		for dealerCard := 2; dealerCard < 11; dealerCard++ {
			decision := rules.MakePlayerDecision(tc.Hand, Card{Value: dealerCard}, 0)
			if hasDefault {
				if defAction != decision {
					t.Fatalf("Unexpected default player decision for `%s` vs `%d` , got %s, expected %s",
						tc.Hand.toString(), dealerCard, decision.ToString(), defAction.ToString())
				}
			} else {
				expectedAction := tc.ExpectedDecisions[dealerCard]
				if expectedAction != decision {
					t.Fatalf("Unexpected player decision for `%s` vs `%d`, got %s, expected %s",
						tc.Hand.toString(), dealerCard, decision.ToString(), expectedAction.ToString())
				}
			}
		}
	}
}

func TestH17SplitHandDecisions(t *testing.T) {
	tests := []struct {
		Hand              Hand
		ExpectedDecisions map[int]PlayerDecision
	}{
		// AA
		{Hand: MakeHand(11, 11), ExpectedDecisions: map[int]PlayerDecision{
			-1: PlayerDecisionSplitAces,
		}},

		// 10/10
		{Hand: MakeHand(10, 10), ExpectedDecisions: map[int]PlayerDecision{
			-1: PlayerDecisionStand,
		}},

		// 9/9
		{Hand: MakeHand(9, 9), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionSplit,
			3:  PlayerDecisionSplit,
			4:  PlayerDecisionSplit,
			5:  PlayerDecisionSplit,
			6:  PlayerDecisionSplit,
			7:  PlayerDecisionHit,
			8:  PlayerDecisionSplit,
			9:  PlayerDecisionSplit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// 8/8
		{Hand: MakeHand(8, 8), ExpectedDecisions: map[int]PlayerDecision{
			-1: PlayerDecisionSplit,
		}},
		// 7/7
		{Hand: MakeHand(7, 7), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionSplit,
			3:  PlayerDecisionSplit,
			4:  PlayerDecisionSplit,
			5:  PlayerDecisionSplit,
			6:  PlayerDecisionSplit,
			7:  PlayerDecisionSplit,
			8:  PlayerDecisionHit,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// 6/6
		{Hand: MakeHand(6, 6), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionSplit,
			3:  PlayerDecisionSplit,
			4:  PlayerDecisionSplit,
			5:  PlayerDecisionSplit,
			6:  PlayerDecisionSplit,
			7:  PlayerDecisionHit,
			8:  PlayerDecisionHit,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// 5/5
		{Hand: MakeHand(5, 5), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionDouble,
			3:  PlayerDecisionDouble,
			4:  PlayerDecisionDouble,
			5:  PlayerDecisionDouble,
			6:  PlayerDecisionDouble,
			7:  PlayerDecisionDouble,
			8:  PlayerDecisionDouble,
			9:  PlayerDecisionDouble,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// 4/4
		{Hand: MakeHand(4, 4), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionHit,
			3:  PlayerDecisionHit,
			4:  PlayerDecisionHit,
			5:  PlayerDecisionSplit,
			6:  PlayerDecisionSplit,
			7:  PlayerDecisionHit,
			8:  PlayerDecisionHit,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// 3/3
		{Hand: MakeHand(3, 3), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionSplit,
			3:  PlayerDecisionSplit,
			4:  PlayerDecisionSplit,
			5:  PlayerDecisionSplit,
			6:  PlayerDecisionSplit,
			7:  PlayerDecisionSplit,
			8:  PlayerDecisionHit,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
		// 2/2
		{Hand: MakeHand(2, 2), ExpectedDecisions: map[int]PlayerDecision{
			2:  PlayerDecisionSplit,
			3:  PlayerDecisionSplit,
			4:  PlayerDecisionSplit,
			5:  PlayerDecisionSplit,
			6:  PlayerDecisionSplit,
			7:  PlayerDecisionSplit,
			8:  PlayerDecisionHit,
			9:  PlayerDecisionHit,
			10: PlayerDecisionHit,
			11: PlayerDecisionHit,
		}},
	}

	rules := MakeTestRules()

	for _, tc := range tests {
		defAction, hasDefault := tc.ExpectedDecisions[-1]
		for dealerCard := 2; dealerCard < 11; dealerCard++ {
			decision := rules.MakePlayerDecision(tc.Hand, Card{Value: dealerCard}, 0)
			split := decision == PlayerDecisionSplitAces || decision == PlayerDecisionSplit
			if hasDefault {
				if defAction != decision {
					t.Fatalf("Unexpected default player decision for `%s` vs `%d` , got %s, expected %s",
						tc.Hand.toString(), dealerCard, decision.ToString(), defAction.ToString())
				}
			} else {
				expectedAction := tc.ExpectedDecisions[dealerCard]
				expectedSplit := decision == PlayerDecisionSplitAces || decision == PlayerDecisionSplit
				if expectedSplit != split {
					t.Fatalf("Unexpected player decision for `%s` vs `%d`, got %s, expected %s",
						tc.Hand.toString(), dealerCard, decision.ToString(), expectedAction.ToString())
				}
			}
		}
	}
}

func TestH17DealerDecisions(t *testing.T) {
	tests := []struct {
		Hand              Hand
		ExpectedDecisions PlayerDecision
	}{
		// Hard hands
		{Hand: MakeHand(10, 9), ExpectedDecisions: PlayerDecisionStand},
		{Hand: MakeHand(10, 7), ExpectedDecisions: PlayerDecisionStand},
		{Hand: MakeHand(10, 6), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(10, 5), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(10, 4), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(10, 3), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(10, 2), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(9, 2), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(8, 2), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(7, 2), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(6, 2), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(5, 2), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(4, 2), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(3, 2), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(2, 2), ExpectedDecisions: PlayerDecisionHit},

		// Soft hands
		{Hand: MakeHand(11, 11), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(11, 10), ExpectedDecisions: PlayerDecisionStand},
		{Hand: MakeHand(11, 9), ExpectedDecisions: PlayerDecisionStand},
		{Hand: MakeHand(11, 8), ExpectedDecisions: PlayerDecisionStand},
		{Hand: MakeHand(11, 7), ExpectedDecisions: PlayerDecisionStand},
		{Hand: MakeHand(11, 6), ExpectedDecisions: PlayerDecisionStand},
		{Hand: MakeHand(11, 5), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(11, 4), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(11, 3), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(11, 2), ExpectedDecisions: PlayerDecisionHit},
		// Random

		{Hand: MakeHand(11, 11, 11, 11), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(11, 11, 7, 2), ExpectedDecisions: PlayerDecisionStand},
		{Hand: MakeHand(5, 5, 3, 2), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(5, 5, 3, 3), ExpectedDecisions: PlayerDecisionHit},
		{Hand: MakeHand(5, 5, 3, 3, 11), ExpectedDecisions: PlayerDecisionStand},
	}

	rules := MakeTestRules().SetDealerHitsSoft17(false)
	for _, tc := range tests {
		for dealerCard := 2; dealerCard < 11; dealerCard++ {
			decision := rules.MakeDealerDecision(tc.Hand)
			if decision != tc.ExpectedDecisions {
				t.Fatalf("Unexpected dealer decision for `%s`, got %s, expected %s",
					tc.Hand.toString(), decision.ToString(), tc.ExpectedDecisions.ToString())
			}
		}
	}
}
