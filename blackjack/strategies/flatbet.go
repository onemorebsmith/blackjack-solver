package strategies

import "github.com/onemorebsmith/blackjack-solver/blackjack/core"

type FlatbetStrategy struct{}

func InitFlatbetStrategy() *FlatbetStrategy {
	return &FlatbetStrategy{}
}

func (strat *FlatbetStrategy) Instance() TrackingStrategy { return strat }

func (strat *FlatbetStrategy) Update(cards ...core.Card) {}

func (strat *FlatbetStrategy) Shuffle() {}

func (strat *FlatbetStrategy) Bid(d core.Deck) BidStrategy {
	return BidStrategy{Hands: 1, Units: 1}
}
