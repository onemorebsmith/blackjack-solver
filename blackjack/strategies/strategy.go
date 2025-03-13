package strategies

import "github.com/onemorebsmith/blackjack-solver/blackjack/core"

type TrackingStrategy interface {
	Instance() TrackingStrategy
	Update(cards ...core.Card)
	Bid(d core.Deck) BidStrategy
	Shuffle()
}
