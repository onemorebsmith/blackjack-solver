package strategies

import "github.com/onemorebsmith/blackjack-solver/blackjack/core"

type TrackingStrategy interface {
	Update(cards ...core.Card)
	Bid(d core.Deck) BidStrategy
	Shuffle()
}
