package blackjack

import (
	"math"
	"math/rand"
	"strings"
)

const SuiteSize = 13
const Suits = 4
const DeckSize = SuiteSize * Suits

type Deck struct {
	Cards    []Card
	idx      int
	Count    int
	deckSize int
}

// Creates `shoe` of 1+ decks, unshuffled initially
func GenerateShoe(decks int) *Deck {
	shoe := GenerateDeck()
	for i := 1; i < decks; i++ {
		additionalDeck := GenerateDeck()
		shoe.Cards = append(shoe.Cards, additionalDeck.Cards...)
	}

	return shoe
}

// Creates a non-suffled full 52 card deck
func GenerateDeck() *Deck {
	all := make([]Card, 0, DeckSize)
	for suiteIdx := SuiteFirst; suiteIdx < SuiteLast; suiteIdx++ {
		for cardIdx := 0; cardIdx < SuiteSize; cardIdx++ {
			card := cards[cardIdx]
			card.Suit = Suit(suiteIdx)
			all = append(all, card)
		}
	}

	return &Deck{
		idx:   0,
		Cards: all,
	}
}

func (d *Deck) ToString() string {
	cardsLeft := len(d.Cards[d.idx:])
	s := make([]string, 0, cardsLeft)
	for i := d.idx; i < cardsLeft; i++ {
		s = append(s, d.Cards[i].ToString())
	}
	return strings.Join(s, " ")
}

func (d *Deck) Shuffle() *Deck {
	d.idx = 0
	d.Count = 0
	d.deckSize = len(d.Cards)
	for i := 0; i < 100000; i++ {
		idxA := rand.Intn(d.deckSize)
		idxB := rand.Intn(d.deckSize)
		if idxA != idxB { // swap
			c := d.Cards[idxA]
			d.Cards[idxA] = d.Cards[idxB]
			d.Cards[idxB] = c
		}
	}

	return d
}

func (d *Deck) TrueCount() int {
	return int(math.Floor(float64(d.Count) / ((float64(d.deckSize - d.idx)) / float64(DeckSize))))
}

func (d *Deck) Deal() Card {
	c := d.Cards[d.idx]
	d.idx++
	d.Count += c.CountValue
	return c
}

func (d *Deck) Remaining() int {
	return len(d.Cards) - d.idx
}
