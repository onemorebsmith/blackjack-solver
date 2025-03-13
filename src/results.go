package blackjack

type GameResults struct {
	Hands            int
	Wins             int
	Losses           int
	Pushes           int
	Blackjacks       int
	EV               float32
	Result           float32
	AvgTC            float32
	HighTC           float32
	LowTC            float32
	EVVariance       float32
	HourlyEVVariance float32
	BidsByTC         map[int]int
	HandAVs          []float32
}

func AggregateResults(results ...GameResults) GameResults {
	aggregated := GameResults{BidsByTC: map[int]int{}}
	for _, r := range results {
		aggregated.EV += r.EV
		aggregated.Hands += r.Hands
		aggregated.Blackjacks += r.Blackjacks
		aggregated.Wins += r.Wins
		aggregated.Losses += r.Losses
		aggregated.Pushes += r.Pushes
		aggregated.AvgTC += r.AvgTC
		aggregated.HighTC += r.HighTC
		aggregated.LowTC += r.LowTC

		for tc, freq := range r.BidsByTC {
			aggregated.BidsByTC[tc] += freq
		}
	}
	return aggregated
}
