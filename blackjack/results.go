package blackjack

type GameResults struct {
	Hands      int
	Wins       int
	Losses     int
	Pushes     int
	Blackjacks int
	EV         float32
	Result     float32
	AvgTc      float32
	TCMap      map[int]int
}

func AggregateResults(results ...GameResults) GameResults {
	aggregated := GameResults{}
	aggregated.TCMap = map[int]int{}
	for _, r := range results {
		aggregated.EV += r.EV
		aggregated.Hands += r.Hands
		aggregated.Blackjacks += r.Blackjacks
		aggregated.Wins += r.Wins
		aggregated.Losses += r.Losses
		aggregated.Pushes += r.Pushes
		aggregated.AvgTc += r.AvgTc
		for count, freq := range r.TCMap {
			aggregated.TCMap[count] += freq
		}
	}
	return aggregated
}
