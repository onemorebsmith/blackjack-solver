package blackjack

type BidStrategy struct {
	Hands int
	Units float32
}

type Bidspread struct {
	spread   map[int]BidStrategy
	maxCount int
	maxBet   BidStrategy
}

func NewBidspread(spread map[int]BidStrategy) *Bidspread {
	maxSpread := 0
	maxBet := BidStrategy{Hands: 1, Units: 1}
	for k, v := range spread {
		if k >= maxSpread {
			maxSpread = k
			maxBet = v
		}
	}
	return &Bidspread{
		spread:   spread,
		maxCount: maxSpread,
		maxBet:   maxBet,
	}
}

func (bs *Bidspread) Bid(d *Deck) BidStrategy {
	count := d.TrueCount()
	if count >= bs.maxCount {
		return bs.maxBet
	}
	if bid, exists := bs.spread[count]; exists {
		return bid
	}
	return BidStrategy{Hands: 1, Units: 1}
}
