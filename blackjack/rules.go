package blackjack

import "log"

type HandResult int

const (
	HandResultPush HandResult = iota
	HandResultWin
	HandResultBlackjack
	HandResultDealerBlackjack
	HandResultInsuranceSave
	HandResultLose
)

func ResultToString(h HandResult) string {
	switch h {
	case HandResultPush:
		return `push`
	case HandResultBlackjack:
		return `blackjack`
	case HandResultLose:
		return `lose`
	case HandResultDealerBlackjack:
		return `dealer blackjack`
	case HandResultWin:
		return `win`
	}
	return `unknown`
}

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
	DealerUpCard  int
	PlayerValue   int
	Soft          bool
	PlayerHits    bool
	PlayerDoubles bool
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
	PlayerDecisionHit
	PlayerDecisionDouble
	PlayerDecisionSplit
	PlayerDecisionSplitAces
)

func (d PlayerDecision) ToString() string {
	switch d {
	case PlayerDecisionHit:
		return `hit`
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

func (rs *BlackjackGameRules) MakeDealerDecision(dealerCards Hand) PlayerDecision {
	value, soft := dealerCards.HandValue()
	if value > 17 {
		return PlayerDecisionStand
	}
	if value == 17 {
		if soft && rs.DealerHitsSoft17 {
			return PlayerDecisionHit
		}
		return PlayerDecisionStand
	}
	return PlayerDecisionHit
}

func (rs *BlackjackGameRules) MakePlayerDecision(playerCards Hand, dealerUpcard Card, splitCounter int) PlayerDecision {
	natural := len(playerCards.Cards) == 2
	playerValue, soft := playerCards.HandValue()
	if playerValue == 21 {
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
		if natural && rule.PlayerDoubles {
			return PlayerDecisionDouble
		} else if rule.PlayerHits {
			return PlayerDecisionHit
		}
	} else {
		log.Printf("\tMissing rule: dealer %d vs player %d, soft %t", dealerUpcard.Value, playerValue, soft)
	}
	return PlayerDecisionStand
}

func InitGame(rules []RuleShorthand, splits []SplitRule) *Ruleset {
	ruleMap := RuleMap{}
	// default rules, stand at every value. These will be overwritten later
	for dealerCard := 2; dealerCard <= 11; dealerCard++ {
		for playerCard := 2; playerCard <= 21; playerCard++ {
			created := Rule{
				DealerUpCard:  dealerCard,
				PlayerValue:   playerCard,
				PlayerHits:    playerCard < 11,
				PlayerDoubles: playerCard == 11,
				Soft:          false,
			}
			ruleMap[created.Hash()] = created
			createdSoft := Rule{
				DealerUpCard:  dealerCard,
				PlayerValue:   playerCard,
				PlayerHits:    playerCard < 11,
				PlayerDoubles: playerCard == 11,
				Soft:          true,
			}
			ruleMap[createdSoft.Hash()] = createdSoft
		}
	}

	for _, shorthand := range rules {
		existingRules := map[int]struct{}{}
		for _, r := range shorthand.PlayerHitsOn {
			created := Rule{
				DealerUpCard: shorthand.DealerCard,
				PlayerValue:  r,
				PlayerHits:   true,
				Soft:         shorthand.Soft,
			}
			ruleMap[created.Hash()] = created
			existingRules[shorthand.DealerCard] = struct{}{}
		}

		for _, r := range shorthand.PlayerDoublesOn {
			created := Rule{
				DealerUpCard:  shorthand.DealerCard,
				PlayerValue:   r,
				PlayerHits:    false,
				PlayerDoubles: true,
				Soft:          shorthand.Soft,
			}
			ruleMap[created.Hash()] = created
			existingRules[shorthand.DealerCard] = struct{}{}
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
