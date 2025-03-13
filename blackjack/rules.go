package blackjack

import (
	"log"

	"github.com/onemorebsmith/blackjack-solver/blackjack/core"
)

type SplitRule struct {
	PlayerCard   int
	DealerUpcard []int // split against these dealer cards
}

func HashSplit(playerCard, dealerCard int) int64 {
	idx := int64(playerCard)
	idx += int64(dealerCard) << 8
	return idx
}

type Rule struct {
	DealerUpCard int
	PlayerValue  int
	Action       PlayerAction
	Soft         bool
}

type RuleShorthand struct {
	DealerCard      int
	PlayerHitsOn    []int
	PlayerDoublesOn []int
	Soft            bool
}

type RuleMap map[int64]Rule
type SplitMap map[int64]struct{}

type Ruleset struct {
	rules RuleMap
	spits SplitMap
}

func (r Rule) Hash() int64 {
	idx := int64(r.DealerUpCard)
	idx += int64(r.PlayerValue) << 8
	if r.Soft {
		idx += int64(1) << 16
	}
	return idx
}

type PlayerDecision int

const (
	PlayerDecisionStand PlayerDecision = iota
	PlayerDecisionNatural21
	PlayerDecisionHit
	PlayerDecisionDouble
	PlayerDecisionSplit
	PlayerDecisionSplitAces
)

func (d PlayerDecision) ToString() string {
	switch d {
	case PlayerDecisionHit:
		return `hit`
	case PlayerDecisionNatural21:
		return `blackjack`
	case PlayerDecisionStand:
		return `stand`
	case PlayerDecisionDouble:
		return `double`
	case PlayerDecisionSplit:
		return `split`
	case PlayerDecisionSplitAces:
		return `split aces`
	}
	return `unknown`
}

func (rs *BlackjackGameRules) MakeDealerDecision(dealerCards core.Hand) PlayerDecision {
	value, soft := dealerCards.HandValue()
	if value > 17 {
		return PlayerDecisionStand
	} else if value == 17 {
		if soft && rs.DealerHitsSoft17 {
			return PlayerDecisionHit
		}
		return PlayerDecisionStand
	}
	return PlayerDecisionHit
}

func (rs *BlackjackGameRules) MakePlayerDecision(playerCards core.Hand, dealerUpcard core.Card, splitCounter int) PlayerDecision {
	natural := len(playerCards.Cards) == 2
	playerValue, soft := playerCards.HandValue()
	if natural && playerValue == 21 { // Natural 21
		return PlayerDecisionNatural21
	}

	if playerValue == 21 { // Natural 21
		return PlayerDecisionStand
	}

	if splitCounter < rs.MaxPlayerSplits {
		if val, isPair := playerCards.IsPair(); isPair {
			hash := HashSplit(val, dealerUpcard.Value)
			if _, exists := rs.playerStrategy.spits[hash]; exists {
				if val == 11 {
					return PlayerDecisionSplitAces
				} else {
					return PlayerDecisionSplit
				}
			}
		}
	}
	rule := Rule{
		PlayerValue:  playerValue,
		DealerUpCard: dealerUpcard.Value,
		Soft:         soft,
	}

	if rule, exists := rs.playerStrategy.rules[rule.Hash()]; exists {
		switch rule.Action {
		case PlayerActionDoubleOrHit:
			if natural {
				return PlayerDecisionDouble
			}
			return PlayerDecisionHit
		case PlayerActionDoubleOrStand:
			if natural {
				return PlayerDecisionDouble
			}
			return PlayerDecisionStand
		case PlayerActionHit:
			return PlayerDecisionHit
		case PlayerActionStand:
			return PlayerDecisionStand
		}
	} else {
		log.Printf("\tMissing rule: dealer %d vs player %d, soft %t", dealerUpcard.Value, playerValue, soft)
	}
	return PlayerDecisionStand
}

func InitGame(rules RulesMap, splits []SplitRule) *Ruleset {
	ruleMap := RuleMap{}
	// default rules, stand at every value. These will be overwritten later
	for dealerCard := 2; dealerCard <= 11; dealerCard++ {
		for playerCard := 2; playerCard <= 21; playerCard++ {
			created := Rule{
				DealerUpCard: dealerCard,
				PlayerValue:  playerCard,
				Action:       PlayerActionStand,
				Soft:         false,
			}
			ruleMap[created.Hash()] = created
			createdSoft := Rule{
				DealerUpCard: dealerCard,
				PlayerValue:  playerCard,
				Action:       PlayerActionStand,
				Soft:         true,
			}
			ruleMap[createdSoft.Hash()] = createdSoft
		}
	}

	for dealerCard, rule := range H17Rules {
		for soft, rules := range rule.Actions {
			for playerTotal, action := range rules {
				created := Rule{
					DealerUpCard: dealerCard,
					PlayerValue:  playerTotal,
					Action:       action,
					Soft:         soft,
				}
				ruleMap[created.Hash()] = created
			}
		}
	}

	splitMap := SplitMap{}
	for _, v := range splits {
		for _, dealerCard := range v.DealerUpcard {
			hash := HashSplit(v.PlayerCard, dealerCard)
			splitMap[hash] = struct{}{}
		}
	}

	return &Ruleset{
		rules: ruleMap,
		spits: splitMap,
	}
}
