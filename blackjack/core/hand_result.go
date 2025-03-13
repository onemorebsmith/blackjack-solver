package core

type OverallHandResult int

type HandResult struct {
	Result OverallHandResult
	AV     float32
}

func MakeHandResult(result OverallHandResult, av float32) HandResult {
	return HandResult{Result: result, AV: av}
}

const (
	HandResultPush OverallHandResult = iota
	HandResultBlackjackPush
	HandResultWin
	HandResultBlackjack
	HandResultDealerBlackjack
	HandResultInsuranceSave
	HandResultLose
)

func (h OverallHandResult) ToString() string {
	switch h {
	case HandResultPush:
		return `push`
	case HandResultBlackjackPush:
		return `blackjack push`
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

func AggregateHandResults(results ...HandResult) float32 {
	netAV := float32(0)
	for _, r := range results {
		netAV += r.AV
	}
	return netAV
}
