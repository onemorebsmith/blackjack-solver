package blackjack

type Bidspread struct {
	spread   map[int]float32
	maxCount int
	maxBet   float32
}

func NewBidspread(spread map[int]float32) *Bidspread {
	maxSpread := 0
	maxBet := float32(0)
	for k, v := range spread {
		if k > maxSpread {
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

func (bs *Bidspread) Bid(d *Deck) float32 {
	count := d.TrueCount()
	if count >= bs.maxCount {
		return bs.maxBet
	}
	if bid, exists := bs.spread[count]; exists {
		return bid
	}
	return float32(1)
}
