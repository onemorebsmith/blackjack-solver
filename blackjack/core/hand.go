package core

import (
	"strings"
)

type Hand struct {
	Cards     []Card
	Doubled   bool
	SplitHand bool
}

func (h Hand) HandValue() (int, bool) {
	val := 0
	soft := false
	aceCount := 0
	for _, v := range h.Cards {
		val += v.Value
		if v.Value == 11 {
			soft = true
			aceCount++
		}
	}

	softAces := aceCount
	for i := 0; i < aceCount; i++ {
		if val > 21 && soft {
			val -= 10
			softAces--
		}
	}

	return val, softAces > 0
}

func (h Hand) IsPair() (int, bool) {
	if len(h.Cards) != 2 {
		return 0, false
	}
	return h.Cards[0].Value, h.Cards[0].Value == h.Cards[1].Value
}

func (h Hand) IsNatural() bool {
	return len(h.Cards) == 2
}

func (h Hand) ToString() string {
	s := make([]string, 0, len(h.Cards))
	for _, v := range h.Cards {
		s = append(s, v.ToString())
	}
	return strings.Join(s, " ")
}
