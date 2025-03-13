package blackjack

import (
	"fmt"
	"testing"

	"github.com/onemorebsmith/blackjack-solver/src/blackjack/core"
)

func MakeTestRules() *BlackjackGameRules {
	return NewBlackjackGameRules(InitGame(H17Rules, H17Splits))
}

func MakeHand(values ...int) core.Hand {
	h := core.Hand{}
	for _, v := range values {
		name := fmt.Sprintf("%d", v)
		if v == 11 {
			name = `A`
		}

		h.Cards = append(h.Cards, core.Card{
			Value: v,
			Name:  name,
			Suit:  core.SuitSpades,
		})
	}
	return h
}

func TestSplits(t *testing.T) {
	decision := MakeTestRules().MakePlayerDecision(MakeHand(11, 11), core.Card{Value: 2}, 0)
	if decision != PlayerDecisionSplitAces {
		t.Fatalf("should have split Aces")
	}
	decision = MakeTestRules().MakePlayerDecision(MakeHand(8, 8), core.Card{Value: 2}, 0)
	if decision != PlayerDecisionSplit {
		t.Fatalf("should have split 10s")
	}
	decision = MakeTestRules().MakePlayerDecision(MakeHand(10, 10), core.Card{Value: 2}, 0)
	if decision == PlayerDecisionSplit {
		t.Fatalf("should not split 10s")
	}
}

func TestHits(t *testing.T) {
	decision := MakeTestRules().MakePlayerDecision(MakeHand(3, 5), core.Card{Value: 7}, 0)
	if decision != PlayerDecisionHit {
		t.Fatalf("should have hit on 8 vs 7")
	}
	decision = MakeTestRules().MakePlayerDecision(MakeHand(10, 2), core.Card{Value: 2}, 0)
	if decision != PlayerDecisionHit {
		t.Fatalf("should have hit on 12 vs 2")
	}
	decision = MakeTestRules().MakePlayerDecision(MakeHand(10, 5), core.Card{Value: 10}, 0)
	if decision != PlayerDecisionHit {
		t.Fatalf("should have hit on 15 vs 10")
	}
	decision = MakeTestRules().MakePlayerDecision(MakeHand(11, 5), core.Card{Value: 10}, 0)
	if decision != PlayerDecisionHit {
		t.Fatalf("should have hit on soft 16 vs 10")
	}
	decision = MakeTestRules().MakePlayerDecision(MakeHand(3, 11), core.Card{Value: 8}, 0)
	if decision != PlayerDecisionHit {
		t.Fatalf("should have hit on soft 14 vs 8")
	}
	decision = MakeTestRules().MakePlayerDecision(MakeHand(4, 5, 2), core.Card{Value: 10}, 0)
	if decision != PlayerDecisionHit {
		t.Fatalf("should have hit on hard 11 vs 10")
	}
}

func TestDoubles(t *testing.T) {
	decision := MakeTestRules().MakePlayerDecision(MakeHand(5, 5), core.Card{Value: 7}, 0)
	if decision != PlayerDecisionDouble {
		t.Fatalf("should have doubled 10 vs 7")
	}
	decision = MakeTestRules().MakePlayerDecision(MakeHand(3, 8), core.Card{Value: 10}, 0)
	if decision != PlayerDecisionDouble {
		t.Fatalf("should have doubled 11 vs 10")
	}
	decision = MakeTestRules().MakePlayerDecision(MakeHand(4, 5), core.Card{Value: 2}, 0)
	if decision == PlayerDecisionDouble {
		t.Fatalf("should have doubled 9 vs 2")
	}
}
