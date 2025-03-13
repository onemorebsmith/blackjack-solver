package strategies

import "github.com/onemorebsmith/blackjack-solver/blackjack/core"

type HighLowCountStrategy struct {
	RunningCount int
	betspred     Bidspread
	Updates      int
	HighTC       float32
	LowTC        float32
	AggregatedTC float32
	BidsByTC     map[int]int
}

func InitHighLow(bs map[int]BidStrategy) *HighLowCountStrategy {
	return &HighLowCountStrategy{
		RunningCount: 0,
		betspred:     *NewBidspread(bs),
		BidsByTC:     map[int]int{},
	}
}

func (strat *HighLowCountStrategy) Instance() TrackingStrategy {
	return &HighLowCountStrategy{
		RunningCount: 0,
		betspred:     strat.betspred,
		BidsByTC:     map[int]int{},
	}
}

func (strat *HighLowCountStrategy) Update(cards ...core.Card) {
	for _, c := range cards {
		strat.Updates++
		switch c.Value {
		case 2, 3, 4, 5, 6:
			strat.RunningCount++
		case 10, 11:
			strat.RunningCount--
		default:
		}
	}
}

func (strat *HighLowCountStrategy) Shuffle() {
	strat.RunningCount = 0
	strat.HighTC = 0
	strat.LowTC = 0
}

func (strat *HighLowCountStrategy) Bid(d core.Deck) BidStrategy {
	est := d.EstimateRemaining()
	tc := float32(strat.RunningCount) / est
	if tc < strat.LowTC {
		strat.LowTC = tc
	} else if tc > strat.HighTC {
		strat.HighTC = tc
	}
	strat.AggregatedTC += tc
	strat.BidsByTC[int(tc)]++
	return strat.betspred.Bid(int(tc))
}
